package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"

	"github.com/golang-jwt/jwt/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type JumpstarterClaims struct {
	jwt.RegisteredClaims
	// corev1.ObjectReference
	Kind       string    `json:"kubernetes.io/kind,omitempty"`
	Namespace  string    `json:"kubernetes.io/namespace,omitempty"`
	Name       string    `json:"kubernetes.io/name,omitempty"`
	UID        types.UID `json:"kubernetes.io/uid,omitempty"`
	APIVersion string    `json:"kubernetes.io/api_version,omitempty"`
}

func KeyFunc(_ *jwt.Token) (interface{}, error) {
	key, ok := os.LookupEnv("CONTROLLER_KEY")
	if !ok {
		return nil, fmt.Errorf("Failed to lookup controller key from env")
	}
	return []byte(key), nil
}

func SignObjectToken(
	issuer string,
	audience []string,
	object metav1.Object,
	scheme *runtime.Scheme,
) (string, error) {
	ro, ok := object.(runtime.Object)
	if !ok {
		return "", fmt.Errorf("%T is not a runtime.Object, cannot call SignObjectToken", object)
	}

	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return "", err
	}

	key, err := KeyFunc(nil)
	if err != nil {
		return "", err
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, JumpstarterClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   issuer,
			Subject:  string(object.GetUID()),
			Audience: audience,
			// ExpiresAt: token are valid for the entire lifetime of the object
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        string(uuid.NewUUID()),
		},
		Kind:       gvk.Kind,
		Namespace:  object.GetNamespace(),
		Name:       object.GetName(),
		UID:        object.GetUID(),
		APIVersion: gvk.GroupVersion().String(),
	}).SignedString(key)
}

type Object[T any] interface {
	client.Object
	*T
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
	return nil, fmt.Errorf("no matching client")
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
	var clients jumpstarterdevv1alpha1.ExporterList
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
	return nil, fmt.Errorf("no matching exporter")
}
