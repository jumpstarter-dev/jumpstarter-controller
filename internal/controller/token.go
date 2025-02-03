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
	"k8s.io/apiserver/pkg/apis/apiserver"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/plugin/pkg/authenticator/token/oidc"
	"k8s.io/utils/ptr"
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

var auth, _ = NewJWTAuthenticator(
	context.Background(),
	OIDCConfig{
		JWT: []apiserver.JWTAuthenticator{{
			Issuer: apiserver.Issuer{
				URL: "https://10.239.206.8:5556/dex",
				CertificateAuthority: `-----BEGIN CERTIFICATE-----
MIIB/DCCAYKgAwIBAgIIcpC2uS+SjEIwCgYIKoZIzj0EAwMwIDEeMBwGA1UEAxMV
bWluaWNhIHJvb3QgY2EgNzI5MGI2MCAXDTI1MDIwMzE5MzMyNVoYDzIxMjUwMjAz
MTkzMzI1WjAgMR4wHAYDVQQDExVtaW5pY2Egcm9vdCBjYSA3MjkwYjYwdjAQBgcq
hkjOPQIBBgUrgQQAIgNiAAQzezKJ4My35HPeoJvvzTjhS2uJMBYrYfrs5csxZjiy
q8ORrHM539XhWlA6sVZODhzcF2KL4mC9xKz/yIrsws+LKsIWNHGGmIPEKFYnHBGw
VBGeARvhpzZP/9frJXAN/8ejgYYwgYMwDgYDVR0PAQH/BAQDAgKEMB0GA1UdJQQW
MBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1Ud
DgQWBBSZRBCUuP3ta2xsfjnWIjvgvz4fojAfBgNVHSMEGDAWgBSZRBCUuP3ta2xs
fjnWIjvgvz4fojAKBggqhkjOPQQDAwNoADBlAjADql5Ks5wh181iUa1ZBnx4XOVe
l0l7I+mwlwJSPmkZHxruWZTx7gQU4tfDCr+UuzUCMQC2aDXRb17cphipK4gzbExv
EDLExjhHAqMPrKDmT0jHIi7Bbos38/1tyZ/IoKjLnv0=
-----END CERTIFICATE-----
`,
				Audiences:           []string{"jumpstarter"},
				AudienceMatchPolicy: "MatchAny",
			},
			ClaimValidationRules: []apiserver.ClaimValidationRule{},
			ClaimMappings: apiserver.ClaimMappings{
				Username: apiserver.PrefixedClaimOrExpression{
					Claim:  "sub",
					Prefix: ptr.To(""),
				},
			},
			UserValidationRules: []apiserver.UserValidationRule{},
		}},
	},
	[]string{"jumpstarter"},
	oidc.AllValidSigningAlgorithms(),
	[]string{},
)

func VerifyOIDCToken(ctx context.Context, token string) (user.Info, error) {
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
	token string,
	issuer string,
	audience string,
	kclient client.Client,
) (*jumpstarterdevv1alpha1.Client, error) {
	// Try verify token as an OIDC token, ignore errors
	if userInfo, err := VerifyOIDCToken(ctx, token); err == nil {
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
	}
	return VerifyObjectToken[jumpstarterdevv1alpha1.Client](
		ctx, token, issuer, audience, kclient,
	)
}

func VerifyExporterObjectToken(
	ctx context.Context,
	token string,
	issuer string,
	audience string,
	kclient client.Client,
) (*jumpstarterdevv1alpha1.Exporter, error) {
	// Try verify token as an OIDC token, ignore errors
	if userInfo, err := VerifyOIDCToken(ctx, token); err == nil {
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
	}
	return VerifyObjectToken[jumpstarterdevv1alpha1.Exporter](
		ctx, token, issuer, audience, kclient,
	)
}

func VerifyObjectToken[T any, PT Object[T]](
	ctx context.Context,
	token string,
	issuer string,
	audience string,
	client client.Client,
) (*T, error) {

	parsed, err := jwt.ParseWithClaims(
		token,
		&JumpstarterClaims{},
		KeyFunc,
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Name,
			jwt.SigningMethodHS384.Name,
			jwt.SigningMethodHS512.Name,
		}),
	)
	if err != nil {
		return nil, err
	} else if claims, ok := parsed.Claims.(*JumpstarterClaims); ok {
		var object T
		err = client.Get(
			ctx,
			types.NamespacedName{
				Namespace: claims.Namespace,
				Name:      claims.Name,
			},
			PT(&object),
		)
		if err != nil {
			return nil, err
		}

		if PT(&object).GetUID() != claims.UID {
			return nil, fmt.Errorf("VerifyObjectToken: UID mismatch")
		}

		return &object, nil
	} else {
		return nil, fmt.Errorf("%T is not a JumpstarterClaims", parsed.Claims)
	}
}
