package authorization

import (
	"context"
	"fmt"

	celgo "github.com/google/cel-go/cel"
	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/cel"
	"k8s.io/apiserver/pkg/cel/environment"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CELAuthorizer struct {
	reader   client.Reader
	prefix   string
	compiler cel.Compiler
}
type Expression struct {
	Expression string
}

func (v *Expression) GetExpression() string {
	return v.Expression
}

func (v *Expression) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.BoolType}
}

func NewCELAuthorizer(reader client.Reader, prefix string) (authorizer.Authorizer, error) {
	env, err := environment.MustBaseEnvSet(
		environment.DefaultCompatibilityVersion(),
		false,
	).Extend(environment.VersionedOptions{
		IntroducedVersion: environment.DefaultCompatibilityVersion(),
		EnvOptions: []celgo.EnvOption{
			celgo.Variable("self", celgo.DynType),
			celgo.Variable("user", celgo.DynType),
			celgo.Variable("prefix", celgo.StringType),
			celgo.Variable("kind", celgo.StringType),
		},
	})
	if err != nil {
		return nil, err
	}

	compiler := cel.NewCompiler(env)

	return &CELAuthorizer{
		reader:   reader,
		prefix:   prefix,
		compiler: compiler,
	}, nil
}

func (b *CELAuthorizer) Authorize(
	ctx context.Context,
	attributes authorizer.Attributes,
) (authorizer.Decision, string, error) {
	var self map[string]interface{}
	var err error

	switch attributes.GetResource() {
	case "Exporter":
		var e jumpstarterdevv1alpha1.Exporter
		if err := b.reader.Get(ctx, client.ObjectKey{
			Namespace: attributes.GetNamespace(),
			Name:      attributes.GetName(),
		}, &e); err != nil {
			return authorizer.DecisionDeny, "failed to get exporter", err
		}
		self, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&e)
		if err != nil {
			return authorizer.DecisionDeny, "failed to serialize exporter", err
		}
	case "Client":
		var c jumpstarterdevv1alpha1.Client
		if err := b.reader.Get(ctx, client.ObjectKey{
			Namespace: attributes.GetNamespace(),
			Name:      attributes.GetName(),
		}, &c); err != nil {
			return authorizer.DecisionDeny, "failed to get client", err
		}
		self, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&c)
		if err != nil {
			return authorizer.DecisionDeny, "failed to serialize client", err
		}
	default:
		return authorizer.DecisionDeny, "invalid object kind", nil
	}

	compiled, err := b.compiler.CompileCELExpression(&Expression{
		Expression: "(has(self.spec.username) ? self.spec.username : prefix + kind.lowerAscii() + ':' + self.metadata.namespace + ':' + self.metadata.name + ':' + self.metadata.uid) == user.username",
	})

	user := attributes.GetUser()
	value, _, err := compiled.Program.Eval(map[string]any{
		"self": self,
		"user": map[string]any{
			"username": user.GetName(),
			"uid":      user.GetUID(),
			"groups":   user.GetGroups(),
			"extra":    user.GetExtra(),
		},
		"prefix": b.prefix,
		"kind":   attributes.GetResource(),
	})
	if err != nil {
		return authorizer.DecisionDeny, "failed to evaluate expression", err
	}

	result, ok := value.Value().(bool)
	if !ok {
		return authorizer.DecisionDeny, "failed to evaluate expression", fmt.Errorf("result type mismatch")
	}

	if result {
		return authorizer.DecisionAllow, "", nil
	} else {
		return authorizer.DecisionDeny, "permission denied", nil
	}
}
