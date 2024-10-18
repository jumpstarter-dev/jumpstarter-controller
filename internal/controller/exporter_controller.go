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

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
)

// ExporterReconciler reconciles a Exporter object
type ExporterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jumpstarter.dev,resources=exporters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jumpstarter.dev,resources=exporters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jumpstarter.dev,resources=exporters/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Exporter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *ExporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	exporter := &jumpstarterdevv1alpha1.Exporter{}
	err := r.Get(ctx, req.NamespacedName, exporter)
	if apierrors.IsNotFound(err) {
		logger.Info("reconcile: Exporter deleted", "exporter", req.NamespacedName)
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		return reconcile.Result{}, nil
	}

	if err != nil {
		logger.Error(err, "reconcile: unable to fetch Exporter")
		return ctrl.Result{}, err
	}

	var leases jumpstarterdevv1alpha1.LeaseList
	err = r.List(
		ctx,
		&leases,
		client.InNamespace(req.Namespace),
		MatchingActiveLeases(),
	)
	if err != nil {
		logger.Error(err, "Error listing leases")
		return ctrl.Result{}, err
	}

	exporter.Status.LeaseRef = nil
	for _, lease := range leases.Items {
		if !lease.Status.Ended && lease.Status.ExporterRef != nil {
			if lease.Status.ExporterRef.Name == exporter.Name {
				exporter.Status.LeaseRef = &corev1.LocalObjectReference{Name: lease.Name}
			}
		}
	}
	if err = r.Status().Update(ctx, exporter); err != nil {
		logger.Error(err, "reconcile: unable to update Exporter with leaseRef", "exporter", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if exporter.Status.Credential == nil {
		logger.Info("reconcile: Exporter has no credentials, creating credentials", "exporter", req.NamespacedName)
		secret, err := r.secretForExporter(exporter)
		if err != nil {
			logger.Error(err, "reconcile: unable to create secret for Exporter")
			return ctrl.Result{}, err
		}
		err = r.Create(ctx, secret)
		if err != nil {
			logger.Error(err, "reconcile: unable to create secret for Exporter", "exporter", req.NamespacedName, "secret", secret.GetName())
			return ctrl.Result{}, err
		}
		exporter.Status.Credential = &corev1.LocalObjectReference{
			Name: secret.Name,
		}
		err = r.Status().Update(ctx, exporter)
		if err != nil {
			logger.Error(err, "reconcile: unable to update Exporter with secret reference", "exporter", req.NamespacedName, "secret", secret.GetName())
			return ctrl.Result{}, err
		}
	}

	endpoint := controllerEndpoint()
	if exporter.Status.Endpoint != endpoint {
		logger.Info("reconcile: Exporter endpoint outdated, updating", "exporter", req.NamespacedName)
		exporter.Status.Endpoint = endpoint
		err = r.Status().Update(ctx, exporter)
		if err != nil {
			logger.Error(err, "reconcile: unable to update Exporter with endpoint", "exporter", req.NamespacedName)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ExporterReconciler) secretForExporter(exporter *jumpstarterdevv1alpha1.Exporter) (*corev1.Secret, error) {
	token, err := SignObjectToken(
		"https://jumpstarter.dev/controller",
		[]string{"https://jumpstarter.dev/controller"},
		exporter,
		r.Scheme,
	)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      exporter.Name + "-exporter",
			Namespace: exporter.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"token": token,
		},
	}
	// enable garbage collection on the created resource
	if err := controllerutil.SetOwnerReference(exporter, secret, r.Scheme); err != nil {
		return nil, fmt.Errorf("secretForExporter, error setting owner reference: %w", err)
	}
	return secret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jumpstarterdevv1alpha1.Exporter{}).
		Owns(&jumpstarterdevv1alpha1.Lease{}).
		Complete(r)
}
