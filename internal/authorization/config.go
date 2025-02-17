package authorization

import (
	"context"
	"fmt"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func LoadAuthorizationConfiguration(
	ctx context.Context,
	scheme *runtime.Scheme,
	configuration []byte,
	reader client.Reader,
	prefix string,
) (authorizer.Authorizer, error) {
	var authorizationConfiguration jumpstarterdevv1alpha1.AuthorizationConfiguration
	if err := runtime.DecodeInto(
		serializer.NewCodecFactory(scheme, serializer.EnableStrict).
			UniversalDecoder(jumpstarterdevv1alpha1.GroupVersion),
		configuration,
		&authorizationConfiguration,
	); err != nil {
		return nil, err
	}

	switch authorizationConfiguration.Type {
	case "Basic":
		return NewBasicAuthorizer(reader, prefix), nil
	case "CEL":
		if authorizationConfiguration.CEL == nil {
			return nil, fmt.Errorf("CEL authorizer configuration missing")
		}
		return NewCELAuthorizer(reader, prefix, authorizationConfiguration.CEL.Expression)
	default:
		return nil, fmt.Errorf("unsupported authorizer type: %s", authorizationConfiguration.Type)
	}
}
