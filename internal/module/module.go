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

package module

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kmmv1beta1 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
	examplecomv1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
)

const (
	moduleSuffix = "module"

	driverServiceAccount = "driver-sample"
)

//go:generate mockgen -source=module.go -package=module -destination=mock_module.go

type Reconciler interface {
	ReconcileModule(ctx context.Context, dc *examplecomv1alpha1.DeviceConfig) error
	SetDesiredModule(m *kmmv1beta1.Module, cr *examplecomv1alpha1.DeviceConfig) error
	DeleteModule(ctx context.Context, dc *examplecomv1alpha1.DeviceConfig) error
}

type moduleReconciler struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewReconciler(c client.Client, s *runtime.Scheme) *moduleReconciler {
	return &moduleReconciler{
		client: c,
		scheme: s,
	}
}

func GetModuleName(cr *examplecomv1alpha1.DeviceConfig) string {
	return fmt.Sprintf("%s-%s", cr.Name, moduleSuffix)
}

func (r *moduleReconciler) ReconcileModule(ctx context.Context, cr *examplecomv1alpha1.DeviceConfig) error {
	logger := log.FromContext(ctx)

	existingModule := &kmmv1beta1.Module{}
	err := r.client.Get(ctx, types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      GetModuleName(cr),
	}, existingModule)
	exists := !apierrors.IsNotFound(err)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	m := &kmmv1beta1.Module{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetModuleName(cr),
			Namespace: cr.ObjectMeta.Namespace,
		},
	}

	if exists {
		m = existingModule
	}

	res, err := controllerutil.CreateOrPatch(ctx, r.client, m, func() error {
		return r.SetDesiredModule(m, cr)
	})

	if err != nil {
		return fmt.Errorf("could not create or patch Module: %v", err)
	}

	logger.Info("Reconciled Module", "resource", m.Name, "result", res)

	return nil
}

func (r *moduleReconciler) DeleteModule(ctx context.Context, cr *examplecomv1alpha1.DeviceConfig) error {
	m := &kmmv1beta1.Module{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetModuleName(cr),
			Namespace: cr.ObjectMeta.Namespace,
		},
	}

	err := r.client.Delete(ctx, m)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Module %s: %w", m.Name, err)
	}

	return nil
}

func (r *moduleReconciler) SetDesiredModule(m *kmmv1beta1.Module, cr *examplecomv1alpha1.DeviceConfig) error {
	if m == nil {
		return errors.New("module cannot be nil")
	}

	ModuleLoader := r.makeModuleLoader(cr)
	selector := cr.GetNodeSelector()

	m.Spec = kmmv1beta1.ModuleSpec{
		ModuleLoader: ModuleLoader,
		Selector:     selector,
	}

	if err := ctrl.SetControllerReference(cr, m, r.scheme); err != nil {
		return err
	}

	return nil
}

func (r *moduleReconciler) makeModuleLoader(cr *examplecomv1alpha1.DeviceConfig) kmmv1beta1.ModuleLoaderSpec {
	moduleLoader := kmmv1beta1.ModuleLoaderSpec{
		Container: kmmv1beta1.ModuleLoaderContainerSpec{
			ImagePullPolicy: corev1.PullAlways,
			KernelMappings:  r.makeKernelMappings(cr),
		},
		ServiceAccountName: driverServiceAccount,
	}

	return moduleLoader
}

func (r *moduleReconciler) makeKernelMappings(cr *examplecomv1alpha1.DeviceConfig) []kmmv1beta1.KernelMapping {
	kernelMappings := []kmmv1beta1.KernelMapping{
		{
			ContainerImage: fmt.Sprintf("%s:%s-${KERNEL_FULL_VERSION}", cr.Spec.DriverImage, cr.Spec.DriverVersion),
			Regexp:         `^.*\.el\d_?\d?\..*$`,
		},
	}

	return kernelMappings
}
