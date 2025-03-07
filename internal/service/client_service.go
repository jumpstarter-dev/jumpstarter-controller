package service

import (
	"context"
	"net"

	pb "github.com/jumpstarter-dev/jumpstarter-controller/internal/protocol/jumpstarter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClientService struct {
	pb.UnimplementedClientServiceServer
	Client client.Client
	Scheme *runtime.Scheme
}

func (c *ClientService) GetExporter(
	ctx context.Context,
	req *pb.GetExporterRequest,
) (*pb.Exporter, error) {
	return &pb.Exporter{
		Name:   req.Name,
		Labels: map[string]string{"test": "dummy"},
	}, nil
}

func (c *ClientService) ListExporters(
	ctx context.Context,
	req *pb.ListExportersRequest,
) (*pb.ListExportersResponse, error) {
	return &pb.ListExportersResponse{
		Exporters: []*pb.Exporter{{
			Name:   "namespaces/default/exporters/dummy",
			Labels: map[string]string{"test": "dummy"},
		}},
		NextPageToken: "",
	}, nil
}

func (c *ClientService) Start(ctx context.Context) error {
	server := grpc.NewServer()

	pb.RegisterClientServiceServer(server, c)

	reflection.Register(server)

	listener, err := net.Listen("tcp", "127.0.0.1:8089")
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	return server.Serve(listener)
}

func (c *ClientService) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(c)
}
