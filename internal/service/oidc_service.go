package service

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
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
	Key    []byte
	Cert   *tls.Certificate
}

type ed25519Key struct {
	id  string
	key ed25519.PrivateKey
}

func (key *ed25519Key) ID() string {
	return key.id
}

func (key *ed25519Key) Algorithm() jose.SignatureAlgorithm {
	return jose.EdDSA
}

func (key *ed25519Key) Use() string {
	return "sig"
}

func (key *ed25519Key) Key() any {
	return key.key.Public()
}

func (s *OIDCService) KeySet(context.Context) ([]op.Key, error) {
	return []op.Key{&ed25519Key{
		id:  "default",
		key: ed25519.NewKeyFromSeed(s.Key),
	}}, nil
}

func (s *OIDCService) Start(ctx context.Context) error {
	r := gin.Default()

	r.GET("/.well-known/openid-configuration", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"issuer":   Issuer,
			"jwks_uri": Issuer + "/jwks",
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
