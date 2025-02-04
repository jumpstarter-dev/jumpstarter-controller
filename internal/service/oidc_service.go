package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Issuer = "https://localhost:8085"
)

// RouterService exposes a gRPC service
type OIDCService struct {
	client.Client
	Scheme *runtime.Scheme
	Key    *ecdsa.PrivateKey
	Cert   *tls.Certificate
}

type ecdsaKey struct {
	id  string
	key *ecdsa.PrivateKey
}

func (key *ecdsaKey) ID() string {
	return key.id
}

func (key *ecdsaKey) Algorithm() jose.SignatureAlgorithm {
	return jose.ES256
}

func (key *ecdsaKey) Use() string {
	return "sig"
}

func (key *ecdsaKey) Key() any {
	return key.key.Public()
}

func (s *OIDCService) KeySet(context.Context) ([]op.Key, error) {
	return []op.Key{&ecdsaKey{
		id:  "default",
		key: s.Key,
	}}, nil
}

func (s *OIDCService) Start(ctx context.Context) error {
	r := gin.Default()

	r.GET("/.well-known/openid-configuration", func(c *gin.Context) {
		op.Discover(c.Writer, &oidc.DiscoveryConfiguration{
			Issuer:  Issuer,
			JwksURI: Issuer + "/jwks",
		})
	})

	r.GET("/jwks", func(c *gin.Context) {
		op.Keys(c.Writer, c.Request, s)
	})

	lis, err := net.Listen("tcp", "127.0.0.1:8085")
	if err != nil {
		return err
	}

	tlslis := tls.NewListener(lis, &tls.Config{
		Certificates: []tls.Certificate{*s.Cert},
	})

	return r.RunListener(tlslis)
}

// SetupWithManager sets up the controller with the Manager.
func (s *OIDCService) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(s)
}
