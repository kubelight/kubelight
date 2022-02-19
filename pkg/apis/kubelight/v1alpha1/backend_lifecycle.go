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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
)

var backendCondSet = apis.NewLivingConditionSet()

// GetGroupVersionKind implements kmeta.OwnerRefable
func (b *Backend) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Backend")
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (b *Backend) GetConditionSet() apis.ConditionSet {
	return backendCondSet
}

// InitializeConditions sets the initial values to the conditions.
func (bs *BackendStatus) InitializeConditions() {
	backendCondSet.Manage(bs).InitializeConditions()
}

func (bs *BackendStatus) MarkBackendUnavailable() {
	backendCondSet.Manage(bs).MarkFalse(
		BackendConditionReady,
		"BackendUnavailable",
		"BackendUnavailable")
}

func (bs *BackendStatus) MarkBackendAvailable() {
	backendCondSet.Manage(bs).MarkTrue(BackendConditionReady)
}
