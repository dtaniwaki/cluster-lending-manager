/*
Copyright 2022 Daisuke Taniwaki..

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterlendingmanagerv1alpha1 "github.com/dtaniwaki/cluster-lending-manager/api/v1alpha1"
)

// LendingConfigReconciler reconciles a LendingConfig object
type LendingConfigReconciler struct {
	client.Client
	Recorder record.EventRecorder
	Cron     *Cron
}

const finalizerName = "clusterlendingmanager.dtaniwaki.github.com/finalizer"

//+kubebuilder:rbac:groups=clusterlendingmanager.dtaniwaki.github.com,resources=lendingconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=clusterlendingmanager.dtaniwaki.github.com,resources=lendingconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=clusterlendingmanager.dtaniwaki.github.com,resources=lendingconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;update;patch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LendingConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *LendingConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the LendingConfig instance.
	logger.Info("Fetch LendingConfig")
	config := &LendingConfig{}
	err := r.Get(ctx, req.NamespacedName, config.ToCompatible())
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Handle deleted resources.
	if !config.ObjectMeta.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(config.ToCompatible(), finalizerName) {
			logger.Info("Clear schedules")
			if err := config.ClearSchedules(ctx, r); err != nil {
				logger.Error(err, "Failed to clear schedules")
			}

			controllerutil.RemoveFinalizer(config.ToCompatible(), finalizerName)
			if err := r.Update(ctx, config.ToCompatible()); err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Set finalizer.
	if !controllerutil.ContainsFinalizer(config.ToCompatible(), finalizerName) {
		logger.Info("Set finalizer")
		config.ObjectMeta.Finalizers = append(config.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(ctx, config.ToCompatible()); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Update the schedules.
	logger.Info("Update schedules")
	if err := config.UpdateSchedules(ctx, r); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LendingConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterlendingmanagerv1alpha1.LendingConfig{}).
		Complete(r)
}
