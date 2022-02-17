/*
Copyright 2019 The Knative Authors

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
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// Pool defined the desired state of the cluster pool that user wants to create
//
// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the desired state of the Pool (from the client).
	Spec PoolSpec `json:"spec"`

	// Status communicates the observed state of the Pool (from the controller).
	// +optional
	Status PoolStatus `json:"status,omitempty"`
}

var (
	// Check that Pool can be validated and defaulted.
	_ apis.Validatable   = (*Pool)(nil)
	_ apis.Defaultable   = (*Pool)(nil)
	_ kmeta.OwnerRefable = (*Pool)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Pool)(nil)
)

const (
	// PoolConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	PoolConditionReady = apis.ConditionReady
)

// PoolSpec holds the desired state of the Pool (from the client).
type PoolSpec struct {
	// List of backends that will support the Pool
	Backends []PoolBackend `json:"backends"`
	// Number of clusters that will be running in the Pool
	Count int `json:"count"`
	// Flavor of Kubernetes or OpenShift to run in the Pool
	Flavor Flavor `json:"flavor"`
	// Configuration of the Pool
	Config PoolConfig `json:"config"`
	// Strategy to support for the Pool
	Strategy string `json:"strategy"`
}

type PoolConfig struct {
	// Number of nodes to be run in one cluster
	Nodes int `json:"nodes"`
}

type PoolBackend struct {
	// Name of the backend to be used for this pool
	Name string `json:"name"`
	// Type of the backend to be used
	Type string `json:"type"`
}

type Flavor struct {
	// Name of the flavor - Kubernetes or OpenShift
	Name string `json:"name"`
	// Kubernetes or OpenShift version to be run
	Version string `json:"version"`
}

// PoolStatus communicates the observed state of the Pool (from the controller).
type PoolStatus struct {
	duckv1.Status `json:",inline"`
}

// PoolList is a list of Pool resources
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pool `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (p *Pool) GetStatus() *duckv1.Status {
	return &p.Status.Status
}
