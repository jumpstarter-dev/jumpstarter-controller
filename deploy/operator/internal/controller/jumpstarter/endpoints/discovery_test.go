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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("DiscoverBaseDomain", func() {
	// Note: These tests require OpenShift CRDs to be available in the test environment.
	// They will be skipped if the CRDs are not present, which is expected in non-OpenShift environments.

	Context("when OpenShift is available", func() {
		BeforeEach(func() {
			// Check if OpenShift CRDs are available
			dns := &unstructured.Unstructured{}
			dns.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "config.openshift.io",
				Version: "v1",
				Kind:    "DNS",
			})
			dns.SetName("cluster")
			dns.Object["spec"] = map[string]interface{}{
				"baseDomain": "test-check.com",
			}

			// Try to create a test DNS object to check if the CRD is available
			err := k8sClient.Create(ctx, dns)
			if err != nil {
				Skip("Skipping OpenShift baseDomain auto-detection tests: OpenShift CRDs not available in test environment")
			}
			// Clean up test object
			_ = k8sClient.Delete(ctx, dns)
		})

		Context("when OpenShift DNS cluster config exists", func() {
			It("should successfully auto-detect baseDomain", func() {
				// Create a mock OpenShift DNS cluster config
				dns := &unstructured.Unstructured{}
				dns.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   "config.openshift.io",
					Version: "v1",
					Kind:    "DNS",
				})
				dns.SetName("cluster")
				dns.Object["spec"] = map[string]interface{}{
					"baseDomain": "example.com",
				}

				// Create the DNS object in the cluster
				err := k8sClient.Create(ctx, dns)
				Expect(err).NotTo(HaveOccurred())

				// Test auto-detection
				detectedBaseDomain, err := DiscoverBaseDomain(ctx, k8sClient, "test-namespace")
				Expect(err).NotTo(HaveOccurred())
				Expect(detectedBaseDomain).To(Equal("test-namespace.apps.example.com"))

				// Cleanup
				err = k8sClient.Delete(ctx, dns)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when OpenShift DNS cluster config has empty baseDomain", func() {
			It("should return an error", func() {
				// Create a mock OpenShift DNS cluster config with empty baseDomain
				dns := &unstructured.Unstructured{}
				dns.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   "config.openshift.io",
					Version: "v1",
					Kind:    "DNS",
				})
				dns.SetName("cluster")
				dns.Object["spec"] = map[string]interface{}{
					"baseDomain": "",
				}

				// Create the DNS object in the cluster
				err := k8sClient.Create(ctx, dns)
				Expect(err).NotTo(HaveOccurred())

				// Test auto-detection with empty baseDomain
				_, err = DiscoverBaseDomain(ctx, k8sClient, "test-namespace")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("baseDomain not found or empty"))

				// Cleanup
				err = k8sClient.Delete(ctx, dns)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("when OpenShift DNS cluster config does not exist", func() {
		It("should return an error", func() {
			// Try to auto-detect when no DNS config exists
			// This test will work even without OpenShift CRDs because it just checks error handling
			_, err := DiscoverBaseDomain(ctx, k8sClient, "test-namespace")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to auto-detect baseDomain"))
		})
	})
})
