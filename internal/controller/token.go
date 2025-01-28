package controller

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
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

func VerifyObjectToken[T any, PT Object[T]](
	ctx context.Context,
	token string,
	issuer string,
	audience string,
	client client.Client,
) (*T, error) {
	provider, err := oidc.NewProvider(ctx, "https://example.com") // TODO: cache provider
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: "jumpstarter",
	})

	claims, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	subject := sha256.Sum256([]byte(claims.Subject))

	var object T
	if err = client.Get(
		ctx,
		types.NamespacedName{
			Namespace: "jumpstarter-lab", // TODO: check namespace
			Name:      hex.EncodeToString(subject[:]),
		},
		PT(&object),
	); err != nil {
		return nil, err
	}

	return &object, nil
}
