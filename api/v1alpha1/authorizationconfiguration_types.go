package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuthorizationConfiguration provides versioned configuration for authorization.
type AuthorizationConfiguration struct {
	metav1.TypeMeta
	Type string            `json:"type"`
	CEL  *CELConfiguration `json:"cel,omitempty"`
}

type CELConfiguration struct {
	Expression string `json:"expression"`
}

func init() {
	SchemeBuilder.Register(&AuthorizationConfiguration{})
}
