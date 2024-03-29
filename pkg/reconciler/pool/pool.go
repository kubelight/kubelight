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
	poolreconciler "github.com/kubelight/kubelight/pkg/client/injection/reconciler/kubelight/v1alpha1/pool"
	kubelightListers "github.com/kubelight/kubelight/pkg/client/listers/kubelight/v1alpha1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
)

// Reconciler implements poolreconciler.Interface for
// Pool resources.
type Reconciler struct {
	backendLister kubelightListers.BackendLister
}

// Check that our Reconciler implements Interface
var _ poolreconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *v1alpha1.Pool) reconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Info("Running reconcile loop for Pool")

	o.Status.MarkServiceAvailable()
	return nil
}
