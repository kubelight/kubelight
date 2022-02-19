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

// Backend is a Knative abstraction that encapsulates the interface by which Knative
// components express a desire to have a particular image cached.
//
// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the desired state of the Backend (from the client).
	Spec BackendSpec `json:"spec"`

	// Status communicates the observed state of the Backend (from the controller).
	// +optional
	Status BackendStatus `json:"status,omitempty"`
}

var (
	// Check that Backend can be validated and defaulted.
	_ apis.Validatable   = (*Backend)(nil)
	_ apis.Defaultable   = (*Backend)(nil)
	_ kmeta.OwnerRefable = (*Backend)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Backend)(nil)
)

// BackendSpec holds the desired state of the Backend (from the client).
type BackendSpec struct {
	// Type of the backend; bare-metal, gke, aws, etc
	Type string `json:"type"`
	// Machines is supported only for bare metal backend
	Machines []Machine `json:"machines,omitempty"`
	// Secret that contains GKE/AWS credentials
	Secret string `json:"secret"`
	// Wake on LAN proxy details
	WakeOnLAN WakeOnLAN `json:"wake-on-lan"`
}

type WakeOnLAN struct {
	// WOL proxy host
	Proxy string `json:"proxy"`
	// SSH credentials for the proxy
	SSHSecret string `json:"ssh-secret"`
	// Command to run for WOL
	Command string `json:"command"`
}

type Machine struct {
	// Hostname or IP of the machine
	Host string `json:"host"`
	// SSH credentials for the machine
	SSHSecret string `json:"ssh-secret"`
	// Should the machine be WOLed if it's down
	WakeOnLan string `json:"wake-on-lan"`
}

const (
	// BackendConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	BackendConditionReady = apis.ConditionReady
)

// BackendStatus communicates the observed state of the Backend (from the controller).
type BackendStatus struct {
	duckv1.Status `json:",inline"`
}

// BackendList is a list of Backend resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Backend `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (b *Backend) GetStatus() *duckv1.Status {
	return &b.Status.Status
}
