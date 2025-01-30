package controller

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
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

type ResourceAccessJumpstarter struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Jumpstarter ResourceAccessJumpstarter `json:"jumpstarter"`
}

type Claims struct {
	Subject        string         `json:"sub"`
	Name           string         `json:"preferred_username"`
	ResourceAccess ResourceAccess `json:"resource_access"`
}

func VerifyToken(ctx context.Context, token string) (*Claims, error) {
	provider, err := oidc.NewProvider(ctx, "http://10.239.206.8:8080/realms/master") // FIXME: cache provider instance
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: "jumpstarter", // FIXME: parameterize client_id
	})

	verified, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	var claims Claims // FIXME: custom claims
	if err := verified.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
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
