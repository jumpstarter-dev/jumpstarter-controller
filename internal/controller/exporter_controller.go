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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

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

	var exporter jumpstarterdevv1alpha1.Exporter
	if err := r.Get(ctx, req.NamespacedName, &exporter); err != nil {
		logger.Error(err, "Reconcile: unable to get exporter", "exporter", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	original := client.MergeFrom(exporter.DeepCopy())

	if err := r.reconcileStatusCredential(ctx, &exporter); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileStatusLeaseRef(ctx, &exporter); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileStatusEndpoint(ctx, &exporter); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Status().Patch(ctx, &exporter, original); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ExporterReconciler) reconcileStatusCredential(
	ctx context.Context,
	exporter *jumpstarterdevv1alpha1.Exporter,
) error {
	logger := log.FromContext(ctx)

	if exporter.Status.Credential == nil {
		logger.Info("reconcileStatusCredential: creating credential for exporter", "exporter", exporter)
		secret, err := r.secretForExporter(exporter)
		if err != nil {
			logger.Error(err, "reconcileStatusCredential: failed to prepare credential for exporter", "exporter", exporter)
			return err
		}
		if err := r.Create(ctx, secret); err != nil {
			logger.Error(err, "reconcileStatusCredential: failed to create credential for exporter", "exporter", exporter)
			return err
		}
		exporter.Status.Credential = &corev1.LocalObjectReference{
			Name: secret.Name,
		}
	}

	return nil
}

func (r *ExporterReconciler) reconcileStatusLeaseRef(
	ctx context.Context,
	exporter *jumpstarterdevv1alpha1.Exporter,
) error {
	logger := log.FromContext(ctx)

	var leases jumpstarterdevv1alpha1.LeaseList
	if err := r.List(
		ctx,
		&leases,
		client.InNamespace(exporter.Namespace),
		MatchingActiveLeases(),
	); err != nil {
		logger.Error(err, "reconcileStatusLeaseRef: failed to list active leases", "exporter", exporter)
		return err
	}

	exporter.Status.LeaseRef = nil
	for _, lease := range leases.Items {
		if !lease.Status.Ended && lease.Status.ExporterRef != nil {
			if lease.Status.ExporterRef.Name == exporter.Name {
				exporter.Status.LeaseRef = &corev1.LocalObjectReference{
					Name: lease.Name,
				}
			}
		}
	}

	return nil
}

// nolint:unparam
func (r *ExporterReconciler) reconcileStatusEndpoint(
	ctx context.Context,
	exporter *jumpstarterdevv1alpha1.Exporter,
) error {
	logger := log.FromContext(ctx)

	endpoint := controllerEndpoint()
	if exporter.Status.Endpoint != endpoint {
		logger.Info("reconcileStatusEndpoint: updating controller endpoint", "exporter", exporter)
		exporter.Status.Endpoint = endpoint
	}

	return nil
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
