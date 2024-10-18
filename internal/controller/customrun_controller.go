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
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	knative "knative.dev/pkg/apis"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CustomRunReconciler reconciles a CustomRun object
type CustomRunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tekton.dev,resources=customruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=customruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=customruns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomRun object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *CustomRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var customrun tektonv1beta1.CustomRun
	if err := r.Get(ctx, req.NamespacedName, &customrun); err != nil {
		logger.Info("reconcile: CustomRun: failed to get", "customrun", req.NamespacedName)
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if (customrun.Spec.CustomRef == nil) && (customrun.Spec.CustomSpec == nil) {
		logger.Info("reconcile: CustomRun: both customRef and customSpec is unset", "customrun", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	if (customrun.Spec.CustomRef != nil) && (customrun.Spec.CustomSpec != nil) {
		logger.Info("reconcile: CustomRun: none of customRef and customSpec is unset", "customrun", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	namespace := customrun.Namespace
	for _, param := range customrun.Spec.Params {
		if param.Name == "namespace" {
			if param.Value.Type == tektonv1beta1.ParamTypeString {
				if param.Value.StringVal != "" {
					namespace = param.Value.StringVal
				}
			}
		}
	}

	if customrun.Spec.CustomRef != nil {
		customRef := customrun.Spec.CustomRef
		if customRef.APIVersion == "jumpstarter.dev/v1alpha1" && customRef.Kind == "Lease" {
			if customRef.Name == "" {
				logger.Info("reconcile: CustomRun: name in customRef is unset", "customrun", req.NamespacedName)
				return ctrl.Result{}, nil
			}

			var lease jumpstarterdevv1alpha1.Lease
			if err := r.Get(
				ctx,
				types.NamespacedName{
					Namespace: namespace,
					Name:      customRef.Name,
				},
				&lease,
			); err != nil {
				logger.Info("reconcile: CustomRun: failed to get Lease referred by customRef", "customrun", req.NamespacedName)
				return ctrl.Result{}, err
			}

			if err := controllerutil.SetOwnerReference(&customrun, &lease, r.Scheme); err != nil {
				logger.Info("reconcile: CustomRun: failed to set Lease ownerReferernces", "customrun", req.NamespacedName)
				return reconcile.Result{}, err
			}

			if err := r.Update(ctx, &lease); err != nil {
				logger.Info("reconcile: CustomRun: unable to update lease", "customrun", req.NamespacedName)
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, r.UpdateStatus(ctx, &customrun, &lease)
		}
	}

	if customrun.Spec.CustomSpec != nil {
		customSpec := customrun.Spec.CustomSpec
		if customSpec.APIVersion == "jumpstarter.dev/v1alpha1" && customSpec.Kind == "Lease" {
			var leaseSpec jumpstarterdevv1alpha1.LeaseSpec
			if err := json.NewDecoder(bytes.NewBuffer(customSpec.Spec.Raw)).Decode(&leaseSpec); err != nil {
				logger.Info("reconcile: CustomRun: unable to decode customSpec", "customrun", req.NamespacedName)
				return reconcile.Result{}, err
			}

			var lease jumpstarterdevv1alpha1.Lease
			err := r.Get(ctx, types.NamespacedName{
				Namespace: namespace,
				Name:      customrun.Name,
			}, &lease)

			if err == nil {
				lease.Spec = leaseSpec

				if err := controllerutil.SetOwnerReference(&customrun, &lease, r.Scheme); err != nil {
					return reconcile.Result{}, err
				}

				if err := r.Update(ctx, &lease); err != nil {
					logger.Info("reconcile: unable to update lease", "customrun", req.NamespacedName)
					return reconcile.Result{}, err
				}

				return reconcile.Result{}, r.UpdateStatus(ctx, &customrun, &lease)
			} else if apierrors.IsNotFound(err) {
				lease.ObjectMeta = metav1.ObjectMeta{
					Namespace: namespace,
					Name:      customrun.Name,
				}
				lease.Spec = leaseSpec

				if err := controllerutil.SetOwnerReference(&customrun, &lease, r.Scheme); err != nil {
					return reconcile.Result{}, err
				}

				if err = r.Create(ctx, &lease); err != nil {
					logger.Info("reconcile: unable to create lease", "customrun", req.NamespacedName)
					return reconcile.Result{}, err
				}

				return reconcile.Result{}, r.UpdateStatus(ctx, &customrun, &lease)
			} else {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

func (r *CustomRunReconciler) UpdateStatus(
	ctx context.Context,
	customrun *tektonv1beta1.CustomRun,
	lease *jumpstarterdevv1alpha1.Lease,
) error {
	// task already completed
	if !customrun.Status.GetCondition(knative.ConditionSucceeded).IsUnknown() {
		return nil
	}

	now := metav1.Now()

	// set startTime if not yet set
	if customrun.Status.StartTime == nil {
		customrun.Status.StartTime = &now
	}

	if meta.IsStatusConditionTrue(
		lease.Status.Conditions,
		string(jumpstarterdevv1alpha1.LeaseConditionTypeReady),
	) {
		customrun.Status.CompletionTime = &now
		customrun.Status.SetCondition(&knative.Condition{
			Type:     knative.ConditionSucceeded,
			Status:   corev1.ConditionTrue,
			Severity: knative.ConditionSeverityInfo,
			LastTransitionTime: knative.VolatileTime{
				Inner: metav1.Now(),
			},
			Reason: "Ready",
		})
		var client jumpstarterdevv1alpha1.Client
		if err := r.Get(
			ctx,
			types.NamespacedName{Namespace: lease.Namespace, Name: lease.Spec.ClientRef.Name},
			&client,
		); err != nil {
			return err
		}
		var secret corev1.Secret
		if err := r.Get(
			ctx,
			types.NamespacedName{Namespace: lease.Namespace, Name: client.Status.Credential.Name},
			&secret,
		); err != nil {
			return err
		}
		token, ok := secret.Data["token"]
		if !ok {
			return fmt.Errorf("token not present in secret")
		}
		customrun.Status.Results = []tektonv1beta1.CustomRunResult{
			{
				Name:  "endpoint",
				Value: client.Status.Endpoint,
			},
			{
				Name:  "token",
				Value: string(token),
			},
			{
				Name:  "lease",
				Value: lease.Name,
			},
		}
	} else {
		if meta.IsStatusConditionTrue(
			lease.Status.Conditions,
			string(jumpstarterdevv1alpha1.LeaseConditionTypeUnsatisfiable),
		) {
			customrun.Status.CompletionTime = &now
			customrun.Status.SetCondition(&knative.Condition{
				Type:     knative.ConditionSucceeded,
				Status:   corev1.ConditionFalse,
				Severity: knative.ConditionSeverityInfo,
				LastTransitionTime: knative.VolatileTime{
					Inner: metav1.Now(),
				},
				Reason: "Unsatisfiable",
			})
		} else {
			customrun.Status.SetCondition(&knative.Condition{
				Type:     knative.ConditionSucceeded,
				Status:   corev1.ConditionUnknown,
				Severity: knative.ConditionSeverityInfo,
				LastTransitionTime: knative.VolatileTime{
					Inner: metav1.Now(),
				},
				Reason: "Pending",
			})
		}
	}

	return r.Status().Update(ctx, customrun)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tektonv1beta1.CustomRun{}).
		Owns(&jumpstarterdevv1alpha1.Lease{}, builder.MatchEveryOwner).
		Complete(r)
}
