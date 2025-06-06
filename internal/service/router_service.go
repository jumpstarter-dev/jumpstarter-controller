/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jumpstarter-dev/jumpstarter-controller/internal/authentication"
	pb "github.com/jumpstarter-dev/jumpstarter-controller/internal/protocol/jumpstarter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RouterService exposes a gRPC service
type RouterService struct {
	pb.UnimplementedRouterServiceServer
	ServerOption grpc.ServerOption
	pending      sync.Map
}

type streamContext struct {
	cancel context.CancelFunc
	stream pb.RouterService_StreamServer
}

func (s *RouterService) authenticate(ctx context.Context) (string, error) {
	token, err := authentication.BearerTokenFromContext(ctx)
	if err != nil {
		return "", err
	}

	parsed, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) { return []byte(os.Getenv("ROUTER_KEY")), nil },
		jwt.WithIssuer("https://jumpstarter.dev/stream"),
		jwt.WithAudience("https://jumpstarter.dev/router"),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Name,
			jwt.SigningMethodHS384.Name,
			jwt.SigningMethodHS512.Name,
		}),
	)

	if err != nil || !parsed.Valid {
		return "", status.Errorf(codes.InvalidArgument, "invalid jwt token")
	}

	return parsed.Claims.GetSubject()
}

func (s *RouterService) Stream(stream pb.RouterService_StreamServer) error {
	ctx := stream.Context()
	logger := log.FromContext(ctx)

	streamName, err := s.authenticate(ctx)
	if err != nil {
		logger.Error(err, "failed to authenticate")
		return err
	}

	logger.Info("streaming", "stream", streamName)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sctx := streamContext{
		cancel: cancel,
		stream: stream,
	}

	actual, loaded := s.pending.LoadOrStore(streamName, sctx)
	if loaded {
		defer actual.(streamContext).cancel()
		logger.Info("forwarding", "stream", streamName)
		return Forward(ctx, stream, actual.(streamContext).stream)
	} else {
		logger.Info("waiting for the other side", "stream", streamName)
		<-ctx.Done()
		return nil
	}
}

func (s *RouterService) Start(ctx context.Context) error {
	log := log.FromContext(ctx)

	dnsnames, ipaddresses, err := endpointToSAN(routerEndpoint())
	if err != nil {
		return err
	}

	cert, err := NewSelfSignedCertificate("jumpstarter router", dnsnames, ipaddresses)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(cert)),
		grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor()),
		s.ServerOption,
	)

	pb.RegisterRouterServiceServer(server, s)

	reflection.Register(server)
	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		return err
	}

	log.Info("Starting grpc router service")
	go func() {
		<-ctx.Done()
		log.Info("Stopping grpc router service")
		server.Stop()
	}()

	return server.Serve(listener)
}

// SetupWithManager sets up the controller with the Manager.
func (s *RouterService) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(s)
}
