package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/jumpstarter-dev/jumpstarter-controller/internal/oidc"
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

func (s *OIDCService) Start(ctx context.Context) error {
	r := gin.Default()

	signer := oidc.NewSigner(s.Key)
	signer.Register(r)

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
