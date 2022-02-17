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

package pool

import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	poolinformer "github.com/kubelight/kubelight/pkg/client/injection/informers/kubelight/v1alpha1/pool"
	poolreconciler "github.com/kubelight/kubelight/pkg/client/injection/reconciler/kubelight/v1alpha1/pool"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	poolInformer := poolinformer.Get(ctx)

	r := &Reconciler{}
	impl := poolreconciler.NewImpl(ctx, r)

	poolInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
