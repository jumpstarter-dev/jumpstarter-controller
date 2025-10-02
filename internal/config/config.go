package config

import (
	"context"
	"fmt"

	"github.com/jumpstarter-dev/jumpstarter-controller/internal/oidc"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func LoadRouterConfiguration(
	ctx context.Context,
	client client.Reader,
	key client.ObjectKey,
) ([]grpc.ServerOption, error) {
	var configmap corev1.ConfigMap
	if err := client.Get(ctx, key, &configmap); err != nil {
		return nil, err
	}

	rawConfig, ok := configmap.Data["config"]
	if !ok {
		return nil, fmt.Errorf("LoadRouterConfiguration: missing config section")
	}

	var config Config
	err := yaml.UnmarshalStrict([]byte(rawConfig), &config)
	if err != nil {
		return nil, err
	}

	serverOptions, err := LoadGrpcConfiguration(config.Grpc)
	if err != nil {
		return nil, err
	}

	return serverOptions, nil
}

func LoadConfiguration(
	ctx context.Context,
	client client.Reader,
	scheme *runtime.Scheme,
	key client.ObjectKey,
	signer *oidc.Signer,
	certificateAuthority string,
) (*ConfigurationResult, error) {
	var configmap corev1.ConfigMap
	if err := client.Get(ctx, key, &configmap); err != nil {
		return nil, err
	}

	rawRouter, ok := configmap.Data["router"]
	if !ok {
		return nil, fmt.Errorf("LoadConfiguration: missing router section")
	}

	var router Router
	if err := yaml.Unmarshal([]byte(rawRouter), &router); err != nil {
		return nil, err
	}

	rawConfig, ok := configmap.Data["config"]
	if !ok {
		return nil, fmt.Errorf("LoadConfiguration: missing config section")
	}

	var config Config
	if err := yaml.UnmarshalStrict([]byte(rawConfig), &config); err != nil {
		return nil, err
	}

	authResult, err := LoadAuthenticationConfiguration(
		ctx,
		scheme,
		config.Authentication,
		signer,
		certificateAuthority,
	)
	if err != nil {
		return nil, err
	}

	serverOptions, err := LoadGrpcConfiguration(config.Grpc)
	if err != nil {
		return nil, err
	}

	// Preprocess configuration values (parse durations, cache expensive operations, etc.)
	if err := config.ExporterOptions.PreprocessConfig(); err != nil {
		return nil, err
	}

	return &ConfigurationResult{
		AuthenticationConfigResult: *authResult,
		Router:                     router,
		ServerOptions:              serverOptions,
		Provisioning:               &config.Provisioning,
		ExporterOptions:            &config.ExporterOptions,
	}, nil
}
