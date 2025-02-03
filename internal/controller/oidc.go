package controller

import (
	"context"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/apis/apiserver"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	tokenunion "k8s.io/apiserver/pkg/authentication/token/union"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/apiserver/plugin/pkg/authenticator/token/oidc"
)

// Reference: https://github.com/kubernetes/kubernetes/blob/v1.32.1/pkg/kubeapiserver/authenticator/config.go#L244
func NewJWTAuthenticator(
	ctx context.Context,
	scheme *runtime.Scheme,
	config jumpstarterdevv1alpha1.AuthenticationConfiguration,
	apiAudiences authenticator.Audiences,
	oidcSigningAlgs []string,
	disallowedIssuers []string,
) (authenticator.Token, error) {
	var jwtAuthenticators []authenticator.Token
	for _, jwtAuthenticator := range config.JWT {
		var oidcCAContent oidc.CAContentProvider
		if len(jwtAuthenticator.Issuer.CertificateAuthority) > 0 {
			var oidcCAError error
			oidcCAContent, oidcCAError = dynamiccertificates.NewStaticCAContent(
				"oidc-authenticator",
				[]byte(jwtAuthenticator.Issuer.CertificateAuthority),
			)
			if oidcCAError != nil {
				return nil, oidcCAError
			}
		}
		var jwtAuthenticatorUnversioned apiserver.JWTAuthenticator
		if err := scheme.Convert(&jwtAuthenticator, &jwtAuthenticatorUnversioned, nil); err != nil {
			return nil, err
		}
		oidcAuth, err := oidc.New(ctx, oidc.Options{
			JWTAuthenticator:     jwtAuthenticatorUnversioned,
			CAContentProvider:    oidcCAContent,
			SupportedSigningAlgs: oidcSigningAlgs,
			DisallowedIssuers:    disallowedIssuers,
		})
		if err != nil {
			return nil, err
		}
		jwtAuthenticators = append(jwtAuthenticators, oidcAuth)
	}
	return authenticator.WrapAudienceAgnosticToken(
		apiAudiences,
		tokenunion.NewFailOnError(jwtAuthenticators...),
	), nil
}
