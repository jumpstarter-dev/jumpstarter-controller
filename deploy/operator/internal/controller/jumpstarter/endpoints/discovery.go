/*
Copyright 2025.

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

package endpoints

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// discoverAPIResource checks if a specific API resource is available in the cluster
// groupVersion should be in the format "group/version" (e.g., "networking.k8s.io/v1", "route.openshift.io/v1")
// kind is the resource kind to look for (e.g., "Ingress", "Route")
func discoverAPIResource(config *rest.Config, groupVersion, kind string) bool {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Log.Error(err, "Failed to create discovery client",
			"groupVersion", groupVersion,
			"kind", kind)
		return false
	}

	apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		// API group not found - resource not available
		return false
	}

	for _, resource := range apiResourceList.APIResources {
		if resource.Kind == kind {
			return true
		}
	}

	return false
}

// DiscoverBaseDomain attempts to auto-detect the baseDomain from OpenShift DNS cluster config
// It returns the detected baseDomain in the format "namespace.apps.baseDomain" for
// OpenShift clusters, or an error if it cannot be determined.
func DiscoverBaseDomain(ctx context.Context, c client.Client, namespace string) (string, error) {
	logger := log.FromContext(ctx)

	// Try to fetch the OpenShift DNS cluster configuration
	dns := &unstructured.Unstructured{}
	dns.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "config.openshift.io",
		Version: "v1",
		Kind:    "DNS",
	})

	err := c.Get(ctx, client.ObjectKey{Name: "cluster"}, dns)
	if err != nil {
		logger.Error(err, "Failed to get OpenShift DNS cluster config - baseDomain cannot be auto-detected")
		return "", fmt.Errorf("failed to auto-detect baseDomain from OpenShift DNS cluster config: %w", err)
	}

	// Extract spec.baseDomain from the DNS object
	spec, found, err := unstructured.NestedMap(dns.Object, "spec")
	if err != nil || !found {
		logger.Error(err, "Failed to get spec from OpenShift DNS cluster config")
		return "", fmt.Errorf("failed to get spec from OpenShift DNS cluster config: spec not found")
	}

	openShiftBaseDomain, found, err := unstructured.NestedString(spec, "baseDomain")
	if err != nil || !found || openShiftBaseDomain == "" {
		logger.Error(err, "Failed to get baseDomain from OpenShift DNS cluster config")
		return "", fmt.Errorf("failed to get baseDomain from OpenShift DNS cluster config: baseDomain not found or empty")
	}

	// Format the baseDomain as "namespace.apps.openShiftBaseDomain"
	// This matches the Helm template behavior when .noNs is false
	detectedBaseDomain := fmt.Sprintf("%s.apps.%s", namespace, openShiftBaseDomain)

	logger.Info("Auto-detected baseDomain from OpenShift DNS cluster config",
		"openShiftBaseDomain", openShiftBaseDomain,
		"detectedBaseDomain", detectedBaseDomain,
		"namespace", namespace)

	return detectedBaseDomain, nil
}
