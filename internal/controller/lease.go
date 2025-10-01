package controller

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
)

func MatchingActiveLeases(prev labels.Selector) labels.Selector {
	// TODO: use field selector once KEP-4358 is stabilized
	// Reference: https://github.com/kubernetes/kubernetes/pull/122717
	requirement, err := labels.NewRequirement(
		string(jumpstarterdevv1alpha1.LeaseLabelEnded),
		selection.DoesNotExist,
		[]string{},
	)

	utilruntime.Must(err)

	return prev.Add(*requirement)
}
