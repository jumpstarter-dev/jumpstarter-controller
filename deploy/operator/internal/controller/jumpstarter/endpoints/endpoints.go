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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/deploy/operator/api/v1alpha1"
	"github.com/jumpstarter-dev/jumpstarter-controller/deploy/operator/internal/utils"
)

// Reconciler provides endpoint reconciliation functionality
type Reconciler struct {
	Client           client.Client
	Scheme           *runtime.Scheme
	IngressAvailable bool
	RouteAvailable   bool
}

// NewReconciler creates a new endpoint reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, config *rest.Config) *Reconciler {
	log := logf.Log.WithName("endpoints-reconciler")

	// Discover API availability at initialization
	ingressAvailable := discoverAPIResource(config, "networking.k8s.io/v1", "Ingress")
	routeAvailable := discoverAPIResource(config, "route.openshift.io/v1", "Route")

	log.Info("API discovery completed",
		"ingressAvailable", ingressAvailable,
		"routeAvailable", routeAvailable)

	return &Reconciler{
		Client:           client,
		Scheme:           scheme,
		IngressAvailable: ingressAvailable,
		RouteAvailable:   routeAvailable,
	}
}

// createOrUpdateService creates or updates a service with proper handling of immutable fields
// and owner references. This is the unified service creation method.
func (r *Reconciler) createOrUpdateService(ctx context.Context, service *corev1.Service, owner metav1.Object) error {
	log := logf.FromContext(ctx)

	existingService := &corev1.Service{}
	existingService.Name = service.Name
	existingService.Namespace = service.Namespace

	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, existingService, func() error {
		// Preserve immutable fields if service already exists
		if existingService.CreationTimestamp.IsZero() {
			// Service is being created, copy all fields from desired service
			existingService.Spec = service.Spec
			existingService.Labels = service.Labels
			existingService.Annotations = service.Annotations
			return controllerutil.SetControllerReference(owner, existingService, r.Scheme)

		} else {
			// Preserve existing NodePorts to prevent "port already allocated" errors
			if service.Spec.Type == corev1.ServiceTypeNodePort || service.Spec.Type == corev1.ServiceTypeLoadBalancer {
				for i := range existingService.Spec.Ports {
					if existingService.Spec.Ports[i].NodePort != 0 && i < len(service.Spec.Ports) {
						service.Spec.Ports[i].NodePort = existingService.Spec.Ports[i].NodePort
					}
				}
			}

			// Update all mutable fields
			if service.Spec.LoadBalancerClass != nil && *service.Spec.LoadBalancerClass != "" {
				existingService.Spec.LoadBalancerClass = service.Spec.LoadBalancerClass
			}
			if service.Spec.ExternalTrafficPolicy != "" {
				existingService.Spec.ExternalTrafficPolicy = service.Spec.ExternalTrafficPolicy
			}

			existingService.Spec.Ports = service.Spec.Ports
			existingService.Spec.Selector = service.Spec.Selector
			existingService.Spec.Type = service.Spec.Type
			existingService.Labels = service.Labels
			existingService.Annotations = service.Annotations
			return controllerutil.SetControllerReference(owner, existingService, r.Scheme)
		}
	})

	if err != nil {
		log.Error(err, "Failed to reconcile service",
			"name", service.Name,
			"namespace", service.Namespace,
			"type", service.Spec.Type)
		return err
	}

	log.Info("Service reconciled",
		"name", service.Name,
		"namespace", service.Namespace,
		"type", service.Spec.Type,
		"selector", service.Spec.Selector,
		"operation", op)

	return nil
}

// ReconcileControllerEndpoint reconciles a controller endpoint service with proper pod selector
// This function creates a separate service for each enabled service type (ClusterIP, NodePort, LoadBalancer)
func (r *Reconciler) ReconcileControllerEndpoint(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort) error {
	// Controller pods have fixed labels: app=jumpstarter-controller
	// We need to create a service with selector matching those labels
	baseLabels := map[string]string{
		"component":  "controller",
		"app":        "jumpstarter-controller",
		"controller": owner.GetName(),
	}

	// Pod selector for controller pods
	podSelector := map[string]string{
		"app": "jumpstarter-controller",
	}

	// Create ingress and route resources
	if err := r.createIngressAndRouteForController(ctx, owner, endpoint, servicePort, baseLabels); err != nil {
		return err
	}

	// Create LoadBalancer service
	if err := r.createLoadBalancerServiceForController(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create NodePort service
	if err := r.createNodePortServiceForController(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create ClusterIP service
	if err := r.createClusterIPServiceForController(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create default service if no service type is enabled
	if err := r.createDefaultServiceForController(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	return nil
}

// createIngressAndRouteForController creates ingress and route resources for controller endpoint
func (r *Reconciler) createIngressAndRouteForController(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, baseLabels map[string]string) error {
	// Ingress resource (uses ClusterIP service)
	if endpoint.Ingress != nil && endpoint.Ingress.Enabled {
		serviceName := servicePort.Name
		if err := r.createIngressForEndpoint(ctx, owner, serviceName, servicePort.Port, endpoint, baseLabels); err != nil {
			return err
		}
	}

	// Route resource (uses ClusterIP service)
	if endpoint.Route != nil && endpoint.Route.Enabled {
		serviceName := servicePort.Name
		if err := r.createRouteForEndpoint(ctx, owner, serviceName, servicePort.Port, endpoint, baseLabels); err != nil {
			return err
		}
	}

	return nil
}

// createLoadBalancerServiceForController creates LoadBalancer service for controller endpoint
func (r *Reconciler) createLoadBalancerServiceForController(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	if endpoint.LoadBalancer != nil && endpoint.LoadBalancer.Enabled {
		return r.createService(ctx, owner, servicePort, "-lb", corev1.ServiceTypeLoadBalancer,
			podSelector, baseLabels, endpoint.LoadBalancer.Annotations, endpoint.LoadBalancer.Labels)
	}
	return nil
}

// createNodePortServiceForController creates NodePort service for controller endpoint
func (r *Reconciler) createNodePortServiceForController(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	if endpoint.NodePort != nil && endpoint.NodePort.Enabled {
		return r.createService(ctx, owner, servicePort, "-np", corev1.ServiceTypeNodePort,
			podSelector, baseLabels, endpoint.NodePort.Annotations, endpoint.NodePort.Labels)
	}
	return nil
}

// createClusterIPServiceForController creates ClusterIP service for controller endpoint
func (r *Reconciler) createClusterIPServiceForController(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	// Create ClusterIP if explicitly enabled OR if Ingress/Route need it
	if (endpoint.ClusterIP != nil && endpoint.ClusterIP.Enabled) ||
		(endpoint.Ingress != nil && endpoint.Ingress.Enabled) ||
		(endpoint.Route != nil && endpoint.Route.Enabled) {
		// Merge annotations and labels from ClusterIP config if present
		var annotations, labels map[string]string
		if endpoint.ClusterIP != nil {
			annotations = endpoint.ClusterIP.Annotations
			labels = endpoint.ClusterIP.Labels
		}
		return r.createService(ctx, owner, servicePort, "", corev1.ServiceTypeClusterIP,
			podSelector, baseLabels, annotations, labels)
	}
	return nil
}

// createDefaultServiceForController creates default ClusterIP service if no service type is enabled
func (r *Reconciler) createDefaultServiceForController(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	// If no service type is explicitly enabled, create a default ClusterIP service
	if (endpoint.LoadBalancer == nil || !endpoint.LoadBalancer.Enabled) &&
		(endpoint.NodePort == nil || !endpoint.NodePort.Enabled) &&
		(endpoint.ClusterIP == nil || !endpoint.ClusterIP.Enabled) &&
		(endpoint.Ingress == nil || !endpoint.Ingress.Enabled) &&
		(endpoint.Route == nil || !endpoint.Route.Enabled) {

		// TODO: Default to Route or Ingress depending of the type of cluster
		return r.createService(ctx, owner, servicePort, "", corev1.ServiceTypeClusterIP,
			podSelector, baseLabels, nil, nil)
	}
	return nil
}

// ReconcileRouterReplicaEndpoint reconciles service, ingress, and route for a specific router replica endpoint
// This function creates a separate service for each enabled service type (ClusterIP, NodePort, LoadBalancer)
func (r *Reconciler) ReconcileRouterReplicaEndpoint(ctx context.Context, owner metav1.Object, replicaIndex int32, endpointIdx int, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort) error {
	// IMPORTANT: The pod selector must match the actual pod labels
	// Router pods have label: app: jumpstarter-router-0 (for replica 0)
	baseAppLabel := fmt.Sprintf("%s-router-%d", owner.GetName(), replicaIndex)

	baseLabels := map[string]string{
		"component":    "router",
		"router":       owner.GetName(),
		"router-index": fmt.Sprintf("%d", replicaIndex),
		"endpoint-idx": fmt.Sprintf("%d", endpointIdx),
	}

	// Pod selector - this MUST match the deployment's pod template labels
	podSelector := map[string]string{
		"app": baseAppLabel, // e.g., "jumpstarter-router-0"
	}

	// Create ingress and route resources
	if err := r.createIngressAndRouteForRouter(ctx, owner, endpoint, servicePort, baseLabels); err != nil {
		return err
	}

	// Create LoadBalancer service
	if err := r.createLoadBalancerServiceForRouter(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create NodePort service
	if err := r.createNodePortServiceForRouter(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create ClusterIP service
	if err := r.createClusterIPServiceForRouter(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	// Create default service if no service type is enabled
	if err := r.createDefaultServiceForRouter(ctx, owner, endpoint, servicePort, podSelector, baseLabels); err != nil {
		return err
	}

	return nil
}

// createIngressAndRouteForRouter creates ingress and route resources for router endpoint
func (r *Reconciler) createIngressAndRouteForRouter(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, baseLabels map[string]string) error {
	// Ingress resource (uses ClusterIP service)
	if endpoint.Ingress != nil && endpoint.Ingress.Enabled {
		serviceName := servicePort.Name
		if err := r.createIngressForEndpoint(ctx, owner, serviceName, servicePort.Port, endpoint, baseLabels); err != nil {
			return err
		}
	}

	// Route resource (uses ClusterIP service)
	if endpoint.Route != nil && endpoint.Route.Enabled {
		serviceName := servicePort.Name
		if err := r.createRouteForEndpoint(ctx, owner, serviceName, servicePort.Port, endpoint, baseLabels); err != nil {
			return err
		}
	}

	return nil
}

// createLoadBalancerServiceForRouter creates LoadBalancer service for router endpoint
func (r *Reconciler) createLoadBalancerServiceForRouter(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	if endpoint.LoadBalancer != nil && endpoint.LoadBalancer.Enabled {
		return r.createService(ctx, owner, servicePort, "-lb", corev1.ServiceTypeLoadBalancer,
			podSelector, baseLabels, endpoint.LoadBalancer.Annotations, endpoint.LoadBalancer.Labels)
	}
	return nil
}

// createNodePortServiceForRouter creates NodePort service for router endpoint
func (r *Reconciler) createNodePortServiceForRouter(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	if endpoint.NodePort != nil && endpoint.NodePort.Enabled {
		return r.createService(ctx, owner, servicePort, "-np", corev1.ServiceTypeNodePort,
			podSelector, baseLabels, endpoint.NodePort.Annotations, endpoint.NodePort.Labels)
	}
	return nil
}

// createClusterIPServiceForRouter creates ClusterIP service for router endpoint
func (r *Reconciler) createClusterIPServiceForRouter(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	// Create ClusterIP if explicitly enabled OR if Ingress/Route need it
	if (endpoint.ClusterIP != nil && endpoint.ClusterIP.Enabled) ||
		(endpoint.Ingress != nil && endpoint.Ingress.Enabled) ||
		(endpoint.Route != nil && endpoint.Route.Enabled) {
		// Merge annotations and labels from ClusterIP config if present
		var annotations, labels map[string]string
		if endpoint.ClusterIP != nil {
			annotations = endpoint.ClusterIP.Annotations
			labels = endpoint.ClusterIP.Labels
		}
		return r.createService(ctx, owner, servicePort, "", corev1.ServiceTypeClusterIP,
			podSelector, baseLabels, annotations, labels)
	}
	return nil
}

// createDefaultServiceForRouter creates default ClusterIP service if no service type is enabled
func (r *Reconciler) createDefaultServiceForRouter(ctx context.Context, owner metav1.Object, endpoint *operatorv1alpha1.Endpoint, servicePort corev1.ServicePort, podSelector map[string]string, baseLabels map[string]string) error {
	// If no service type is explicitly enabled, create a default ClusterIP service
	if (endpoint.LoadBalancer == nil || !endpoint.LoadBalancer.Enabled) &&
		(endpoint.NodePort == nil || !endpoint.NodePort.Enabled) &&
		(endpoint.ClusterIP == nil || !endpoint.ClusterIP.Enabled) &&
		(endpoint.Ingress == nil || !endpoint.Ingress.Enabled) &&
		(endpoint.Route == nil || !endpoint.Route.Enabled) {
		return r.createService(ctx, owner, servicePort, "", corev1.ServiceTypeClusterIP,
			podSelector, baseLabels, nil, nil)
	}
	return nil
}

// createService creates or updates a single service with the specified type and suffix
// This is the unified service creation method that uses createOrUpdateService internally
func (r *Reconciler) createService(ctx context.Context, owner metav1.Object, servicePort corev1.ServicePort,
	nameSuffix string, serviceType corev1.ServiceType, podSelector map[string]string,
	baseLabels map[string]string, annotations map[string]string, extraLabels map[string]string) error {

	// Build service name with suffix to avoid conflicts
	serviceName := servicePort.Name + nameSuffix

	// Merge labels (extra labels take precedence)
	serviceLabels := utils.MergeMaps(baseLabels, extraLabels)

	// Ensure annotations map is initialized
	if annotations == nil {
		annotations = make(map[string]string)
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceName,
			Namespace:   owner.GetNamespace(),
			Labels:      serviceLabels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector: podSelector, // Use the provided pod selector map
			Ports:    []corev1.ServicePort{servicePort},
			Type:     serviceType,
		},
	}

	return r.createOrUpdateService(ctx, service, owner)
}
