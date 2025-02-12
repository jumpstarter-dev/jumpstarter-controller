package service

import (
	"context"
	"net/http"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/jumpstarter-dev/jumpstarter-controller/internal/protocol/jumpstarter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ctrl "sigs.k8s.io/controller-runtime"
)

type GatewayService struct{}

func (g *GatewayService) Start(ctx context.Context) error {
	mux := gwruntime.NewServeMux()

	pb.RegisterClientServiceHandlerFromEndpoint(
		ctx,
		mux,
		"127.0.0.1:8089",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	return http.ListenAndServe(":8088", mux)
}

func (g *GatewayService) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(g)
}
