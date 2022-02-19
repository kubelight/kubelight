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
	"github.com/kubelight/kubelight/pkg/apis/kubelight/v1alpha1"
	backendinformers "github.com/kubelight/kubelight/pkg/client/injection/informers/kubelight/v1alpha1/backend"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/logging"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	poolinformers "github.com/kubelight/kubelight/pkg/client/injection/informers/kubelight/v1alpha1/pool"
	poolreconciler "github.com/kubelight/kubelight/pkg/client/injection/reconciler/kubelight/v1alpha1/pool"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	// Obtain an informer to both the main and child resources. These will be started by
	// the injection framework automatically. They'll keep a cached representation of the
	// cluster's state of the respective resource at all times.
	logger := logging.FromContext(ctx)
	logger.Info("Initializing Pool controller...")
	poolInformer := poolinformers.Get(ctx)
	backendInformer := backendinformers.Get(ctx)

	r := &Reconciler{
		// We need to watch backends for a given pool
		backendLister: backendInformer.Lister(),
	}
	impl := poolreconciler.NewImpl(ctx, r)

	// Listen for events on the main resource and enqueue themselves.
	poolInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Listen for events on the child resources and enqueue the owner of them.
	backendInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&v1alpha1.Pool{}),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
