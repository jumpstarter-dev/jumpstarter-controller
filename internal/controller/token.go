package controller

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"

	"github.com/golang-jwt/jwt/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SignExporterToken(
	issuer string,
	audience []string,
	exporter *jumpstarterdevv1alpha1.Exporter,
	scheme *runtime.Scheme,
	key interface{},
) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strings.TrimPrefix(*exporter.Spec.OIDCSubject, "internal:"),
		Audience:  audience,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
	}).SignedString(key)
}

func SignClientToken(
	issuer string,
	audience []string,
	client *jumpstarterdevv1alpha1.Client,
	scheme *runtime.Scheme,
	key interface{},
) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strings.TrimPrefix(*client.Spec.OIDCSubject, "internal:"),
		Audience:  audience,
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        string(uuid.NewUUID()),
	}).SignedString(key)
}

type OIDCClaims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
}

func VerifyOIDCToken(ctx context.Context, auth authenticator.Token, token string) (user.Info, error) {
	resp, ok, err := auth.AuthenticateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("failed to authenticate token")
	}

	return resp.User, nil
}

func VerifyClientObjectToken(
	ctx context.Context,
	auth authenticator.Token,
	token string,
	issuer string,
	audience string,
	kclient client.Client,
) (*jumpstarterdevv1alpha1.Client, error) {
	userInfo, err := VerifyOIDCToken(ctx, auth, token)
	if err != nil {
		return nil, err
	}
	var clients jumpstarterdevv1alpha1.ClientList
	if err = kclient.List(ctx, &clients); err != nil {
		return nil, err
	}
	for _, c := range clients.Items {
		if true &&
			c.Spec.OIDCSubject != nil &&
			*c.Spec.OIDCSubject == userInfo.GetName() {
			return &c, nil
		}
	}

	name, ok := userInfo.GetExtra()["jumpstarter.dev/name"]
	if !ok || len(name) != 1 {
		return nil, fmt.Errorf("no matching exporter")
	}

	var client = jumpstarterdevv1alpha1.Client{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: os.Getenv("NAMESPACE"), // FIXME: read namespace from claim
			Name:      name[0],
		},
		Spec: jumpstarterdevv1alpha1.ClientSpec{
			OIDCSubject: ptr.To(userInfo.GetName()),
		},
	}
	if err := kclient.Create(ctx, &client); err != nil {
		return nil, err
	}

	return &client, nil
}

func VerifyExporterObjectToken(
	ctx context.Context,
	auth authenticator.Token,
	token string,
	issuer string,
	audience string,
	kclient client.Client,
) (*jumpstarterdevv1alpha1.Exporter, error) {
	userInfo, err := VerifyOIDCToken(ctx, auth, token)
	if err != nil {
		return nil, err
	}
	var exporters jumpstarterdevv1alpha1.ExporterList
	if err = kclient.List(ctx, &exporters); err != nil {
		return nil, err
	}
	for _, e := range exporters.Items {
		if true &&
			e.Spec.OIDCSubject != nil &&
			*e.Spec.OIDCSubject == userInfo.GetName() {
			return &e, nil
		}
	}

	name, ok := userInfo.GetExtra()["jumpstarter.dev/name"]
	if !ok || len(name) != 1 {
		return nil, fmt.Errorf("no matching exporter")
	}

	var exporter = jumpstarterdevv1alpha1.Exporter{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: os.Getenv("NAMESPACE"), // FIXME: read namespace from claim
			Name:      name[0],
		},
		Spec: jumpstarterdevv1alpha1.ExporterSpec{
			OIDCSubject: ptr.To(userInfo.GetName()),
		},
	}
	if err := kclient.Create(ctx, &exporter); err != nil {
		return nil, err
	}

	return &exporter, nil
}
