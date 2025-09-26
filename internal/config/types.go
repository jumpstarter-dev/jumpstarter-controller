package config

import (
	"fmt"
	"time"

	apiserverv1beta1 "k8s.io/apiserver/pkg/apis/apiserver/v1beta1"
)

type Config struct {
	Authentication  Authentication  `json:"authentication"`
	Provisioning    Provisioning    `json:"provisioning"`
	Grpc            Grpc            `json:"grpc"`
	ExporterOptions ExporterOptions `json:"exporterOptions"`
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

type ExporterOptions struct {
	OfflineTimeout    string        `json:"offlineTimeout,omitempty"` // How long to wait before marking the exporter as offline
	offlineTimeoutDur time.Duration // Pre-calculated duration, set during LoadConfiguration
}

// PreprocessConfig parses and caches configuration values that require processing
// This method should be called once during configuration loading to pre-calculate
// expensive operations and cache the results for efficient retrieval
func (e *ExporterOptions) PreprocessConfig() error {
	// Parse and cache the offline timeout duration
	if e.OfflineTimeout == "" {
		e.offlineTimeoutDur = 3 * time.Minute // Default fallback
	} else {
		duration, err := time.ParseDuration(e.OfflineTimeout)
		if err != nil {
			return fmt.Errorf("PreprocessConfig: failed to parse exporter offline timeout: %w", err)
		} else {
			e.offlineTimeoutDur = duration
		}
	}

	// Future configuration parsing can be added here
	return nil
}

// GetOfflineTimeout returns the pre-calculated offline timeout duration
func (e *ExporterOptions) GetOfflineTimeout() time.Duration {
	return e.offlineTimeoutDur
}

type Router map[string]RouterEntry

type RouterEntry struct {
	Endpoint string            `json:"endpoint"`
	Labels   map[string]string `json:"labels"`
}
