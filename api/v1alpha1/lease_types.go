/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LeaseSpec defines the desired state of Lease
type LeaseSpec struct {
	BeginTime metav1.Time `json:"beginTime"`
	EndTime   metav1.Time `json:"endTime"`

	// The client that is requesting the lease
	ClientName string `json:"clientName"`

	// The selector for the exporter to be used
	Selector metav1.LabelSelector `json:"selector"`
}

// LeaseStatus defines the observed state of Lease
type LeaseStatus struct {
	// The minutes that the lease is going to last
	Minutes int `json:"minutes"`
	// If the lease has been acquired an exporter name is assigned
	// and then and then it can be used, it will be empty while still pending
	ExporterName string             `json:"exporterName,omitempty"`
	Ended        bool               `json:"ended"`
	Conditions   []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Lease is the Schema for the exporters API
type Lease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LeaseSpec   `json:"spec,omitempty"`
	Status LeaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LeaseList contains a list of Lease
type LeaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Lease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Lease{}, &LeaseList{})
}