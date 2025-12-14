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
	"fmt"

	operatorv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/deploy/operator/api/v1alpha1"
)

// ApplyEndpointDefaults generates default endpoints for a JumpstarterSpec
// based on the baseDomain and cluster capabilities (Route vs Ingress availability).
func ApplyEndpointDefaults(spec *operatorv1alpha1.JumpstarterSpec, routeAvailable, ingressAvailable bool) {
	// Skip endpoint generation if no baseDomain is set
	if spec.BaseDomain == "" {
		return
	}

	// Generate default controller gRPC endpoint if none specified
	if len(spec.Controller.GRPC.Endpoints) == 0 {
		endpoint := operatorv1alpha1.Endpoint{
			Address: fmt.Sprintf("grpc.%s", spec.BaseDomain),
		}
		// Auto-select Route or Ingress based on cluster capabilities
		if routeAvailable {
			endpoint.Route = &operatorv1alpha1.RouteConfig{Enabled: true}
		} else if ingressAvailable {
			endpoint.Ingress = &operatorv1alpha1.IngressConfig{Enabled: true}
		}
		spec.Controller.GRPC.Endpoints = []operatorv1alpha1.Endpoint{endpoint}
	}

	// Generate default router gRPC endpoints if none specified
	if len(spec.Routers.GRPC.Endpoints) == 0 {
		endpoint := operatorv1alpha1.Endpoint{
			// Use $(replica) placeholder for per-replica addresses
			Address: fmt.Sprintf("router-$(replica).%s", spec.BaseDomain),
		}
		// Auto-select Route or Ingress based on cluster capabilities
		if routeAvailable {
			endpoint.Route = &operatorv1alpha1.RouteConfig{Enabled: true}
		} else if ingressAvailable {
			endpoint.Ingress = &operatorv1alpha1.IngressConfig{Enabled: true}
		}
		spec.Routers.GRPC.Endpoints = []operatorv1alpha1.Endpoint{endpoint}
	}
}
