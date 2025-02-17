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
	reader  client.Reader
	prefix  string
	program celgo.Program
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

func NewCELAuthorizer(reader client.Reader, prefix string, expression string) (authorizer.Authorizer, error) {
	env, err := environment.MustBaseEnvSet(
		environment.DefaultCompatibilityVersion(),
		false,
	).Extend(environment.VersionedOptions{
		IntroducedVersion: environment.DefaultCompatibilityVersion(),
		EnvOptions: []celgo.EnvOption{
			celgo.Variable("kind", celgo.StringType),
			celgo.Variable("self", celgo.DynType),
			celgo.Variable("user", celgo.DynType),
		},
	})
	if err != nil {
		return nil, err
	}

	compiler := cel.NewCompiler(env)

	compiled, err := compiler.CompileCELExpression(&Expression{
		Expression: expression,
	})
	if err != nil {
		return nil, err
	}

	return &CELAuthorizer{
		reader:  reader,
		prefix:  prefix,
		program: compiled.Program,
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
		self["spec"].(map[string]any)["username"] = e.Username(b.prefix)
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
		self["spec"].(map[string]any)["username"] = c.Username(b.prefix)
	default:
		return authorizer.DecisionDeny, "invalid object kind", nil
	}

	user := attributes.GetUser()
	value, _, err := b.program.Eval(map[string]any{
		"kind": attributes.GetResource(),
		"self": self,
		"user": map[string]any{
			"username": user.GetName(),
			"uid":      user.GetUID(),
			"groups":   user.GetGroups(),
			"extra":    user.GetExtra(),
		},
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
