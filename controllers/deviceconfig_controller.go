/*
Copyright 2022.

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
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kmmv1beta1 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
	examplecomv1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
	"github.com/mresvanis/he-sample-operator/internal/conditions"
	"github.com/mresvanis/he-sample-operator/internal/finalizers"
	"github.com/mresvanis/he-sample-operator/internal/metrics"
	"github.com/mresvanis/he-sample-operator/internal/module"
	"github.com/mresvanis/he-sample-operator/internal/nodeselector"
)

// DeviceConfigReconciler reconciles a DeviceConfig object
type DeviceConfigReconciler struct {
	client.Client

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	mr module.Reconciler

	fu finalizers.Updater

	cu conditions.Updater

	nsv nodeselector.Validator
}

func NewDeviceConfigReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	recorder record.EventRecorder,
	mr module.Reconciler,
	fu finalizers.Updater,
	cu conditions.Updater,
	nsv nodeselector.Validator,
) *DeviceConfigReconciler {

	return &DeviceConfigReconciler{
		Client:   client,
		Scheme:   scheme,
		Recorder: recorder,
		mr:       mr,
		fu:       fu,
		cu:       cu,
		nsv:      nsv,
	}
}

//+kubebuilder:rbac:groups=example.com,resources=deviceconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=deviceconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=deviceconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups="kmm.sigs.k8s.io",resources=modules,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *DeviceConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	deviceConfig := &examplecomv1alpha1.DeviceConfig{}
	err := r.Get(ctx, req.NamespacedName, deviceConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			metrics.ReconciliationFailed.WithLabelValues(req.NamespacedName.Name).Set(0)
			logger.Info("DeviceConfig resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get DeviceConfig", "resource", deviceConfig.Name)
		return ctrl.Result{}, err
	}

	if !deviceConfig.ObjectMeta.DeletionTimestamp.IsZero() {
		metrics.ReconciliationFailed.WithLabelValues(deviceConfig.Name).Set(0)

		if r.fu.ContainsDeletionFinalizer(deviceConfig) {
			if err := r.mr.DeleteModule(ctx, deviceConfig); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to delete DeviceConfig resources: %w", err)
			}
			if err := r.fu.RemoveDeletionFinalizer(ctx, deviceConfig); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.nsv.CheckDeviceConfigForConflictingNodeSelector(ctx, deviceConfig); err != nil {
		logger.Error(err, "Failed to validate DeviceConfig", "resource", deviceConfig.Name)
		r.Recorder.Event(
			deviceConfig,
			v1.EventTypeWarning,
			"Error",
			"Conflicting DeviceConfig NodeSelectors found. Please add or update this DeviceConfig's NodeSelector accordingly.",
		)
		metrics.ReconciliationFailed.WithLabelValues(deviceConfig.Name).Set(1)
		return ctrl.Result{}, nil
	}

	if !r.fu.ContainsDeletionFinalizer(deviceConfig) {
		if err := r.fu.AddDeletionFinalizer(ctx, deviceConfig); err != nil {
			return ctrl.Result{}, err
		}
	}

	if err := r.mr.ReconcileModule(ctx, deviceConfig); err != nil {
		if cerr := r.cu.SetConditionsErrored(ctx, deviceConfig, conditions.ReasonModuleFailed, err.Error()); cerr != nil {
			err = fmt.Errorf("%s: %w", err.Error(), cerr)
		}
		metrics.ReconciliationFailed.WithLabelValues(deviceConfig.Name).Set(1)
		return ctrl.Result{}, err
	}

	metrics.ReconciliationFailed.WithLabelValues(deviceConfig.Name).Set(0)

	r.Recorder.Event(
		deviceConfig,
		v1.EventTypeNormal,
		"Reconciled",
		fmt.Sprintf("Succesfully reconciled DeviceConfig %s/%s", deviceConfig.Namespace, deviceConfig.Name),
	)

	return ctrl.Result{}, r.cu.SetConditionsReady(ctx, deviceConfig, "Reconciled", "All resources have been successfully reconciled")
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeviceConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("device-config").
		For(&examplecomv1alpha1.DeviceConfig{}).
		Owns(&kmmv1beta1.Module{}).
		Complete(r)
}
