package config

import (
	apiserverv1beta1 "k8s.io/apiserver/pkg/apis/apiserver/v1beta1"
)

type Config struct {
	Authentication Authentication `json:"authentication"`
	Provisioning   Provisioning   `json:"provisioning"`
	Grpc           Grpc           `json:"grpc"`
}

type Authentication struct {
	Internal Internal                            `json:"internal"`
	JWT      []apiserverv1beta1.JWTAuthenticator `json:"jwt"`
}

type Provisioning struct {
	Enabled bool `json:"enabled"`
}

type Internal struct {
	Prefix string `json:"prefix"`
}

type Grpc struct {
	Keepalive Keepalive `json:"keepalive"`
}

type Keepalive struct {
	// EnforcementPolicy parameters
	MinTime             string `json:"minTime"`
	PermitWithoutStream bool   `json:"permitWithoutStream"`

	// ServerParameters for connection timeout control
	Timeout               string `json:"timeout,omitempty"`               // How long to wait for ping response before closing
	MaxConnectionIdle     string `json:"maxConnectionIdle,omitempty"`     // Max idle time before closing
	MaxConnectionAge      string `json:"maxConnectionAge,omitempty"`      // Max connection lifetime
	MaxConnectionAgeGrace string `json:"maxConnectionAgeGrace,omitempty"` // Grace period after max age
	Time                  string `json:"time,omitempty"`                  // How often server sends pings
}

type Router map[string]RouterEntry

type RouterEntry struct {
	Endpoint string            `json:"endpoint"`
	Labels   map[string]string `json:"labels"`
}
