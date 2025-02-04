package oidc

import (
	"context"
	"fmt"
	"os"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ClientSubject(c *jumpstarterdevv1alpha1.Client, prefix string) string {
	if c.Spec.OIDCSubject == nil {
		return prefix + c.Name
	} else {
		return *c.Spec.OIDCSubject
	}
}

func ExporterSubject(e *jumpstarterdevv1alpha1.Exporter, prefix string) string {
	if e.Spec.OIDCSubject == nil {
		return prefix + e.Name
	} else {
		return *e.Spec.OIDCSubject
	}
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
	prefix string,
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
		if ClientSubject(&c, prefix) == userInfo.GetName() {
			return &c, nil
		}
	}

	name, ok := userInfo.GetExtra()["jumpstarter.dev/name"]
	if !ok || len(name) != 1 {
		return nil, fmt.Errorf("no matching client")
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
	prefix string,
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
		if ExporterSubject(&e, prefix) == userInfo.GetName() {
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
