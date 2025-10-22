package service

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/jumpstarter-dev/jumpstarter-controller/internal/oidc"
	ctrl "sigs.k8s.io/controller-runtime"
)

// RouterService exposes a gRPC service
type OIDCService struct {
	Signer *oidc.Signer
	Cert   *tls.Certificate
}

func (s *OIDCService) Start(ctx context.Context) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	s.Signer.Register(r)

	lis, err := net.Listen("tcp", "127.0.0.1:8085")
	if err != nil {
		return err
	}

	tlslis := tls.NewListener(lis, &tls.Config{
		Certificates: []tls.Certificate{*s.Cert},
	})

	return r.RunListener(tlslis)
}

func (s *OIDCService) NeedLeaderElection() bool {
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (s *OIDCService) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(s)
}
