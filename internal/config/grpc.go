package config

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const defaultGrpcTimeout = 180 * time.Second

func LoadGrpcConfiguration(config Grpc) ([]grpc.ServerOption, error) {
	var serverOptions []grpc.ServerOption

	// Parse EnforcementPolicy parameters
	minTime, err := time.ParseDuration(config.Keepalive.MinTime)
	if err != nil {
		return nil, err
	}

	serverOptions = append(serverOptions, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             minTime,
		PermitWithoutStream: config.Keepalive.PermitWithoutStream,
	}))

	// Parse ServerParameters for connection timeout control
	serverParams := keepalive.ServerParameters{}

	// Timeout: How long to wait for ping response before closing connection
	if config.Keepalive.Timeout != "" {
		timeout, err := time.ParseDuration(config.Keepalive.Timeout)
		if err != nil {
			return nil, fmt.Errorf("LoadGrpcConfiguration: failed to parse Keepalive Timeout: %w", err)
		}
		serverParams.Timeout = timeout
	} else {
		serverParams.Timeout = defaultGrpcTimeout
	}

	// MaxConnectionIdle: Max idle time before closing connection
	if config.Keepalive.MaxConnectionIdle != "" {
		maxIdle, err := time.ParseDuration(config.Keepalive.MaxConnectionIdle)
		if err != nil {
			return nil, err
		}
		serverParams.MaxConnectionIdle = maxIdle
	}

	// MaxConnectionAge: Max connection lifetime
	if config.Keepalive.MaxConnectionAge != "" {
		maxAge, err := time.ParseDuration(config.Keepalive.MaxConnectionAge)
		if err != nil {
			return nil, err
		}
		serverParams.MaxConnectionAge = maxAge
	}

	// MaxConnectionAgeGrace: Grace period after max age
	if config.Keepalive.MaxConnectionAgeGrace != "" {
		maxAgeGrace, err := time.ParseDuration(config.Keepalive.MaxConnectionAgeGrace)
		if err != nil {
			return nil, err
		}
		serverParams.MaxConnectionAgeGrace = maxAgeGrace
	}

	// Time: How often server sends pings
	if config.Keepalive.Time != "" {
		time, err := time.ParseDuration(config.Keepalive.Time)
		if err != nil {
			return nil, err
		}
		serverParams.Time = time
	}

	// Only add ServerParameters if at least one parameter is set
	if serverParams != (keepalive.ServerParameters{}) {
		serverOptions = append(serverOptions, grpc.KeepaliveParams(serverParams))
	}

	return serverOptions, nil
}
