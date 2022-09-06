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
	"errors"
	"fmt"
	"time"

	gomock "github.com/golang/mock/gomock"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	record "k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kmmv1beta1 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
	examplecomv1alpha1 "github.com/mresvanis/he-sample-operator/api/v1alpha1"
	"github.com/mresvanis/he-sample-operator/internal/client"
	"github.com/mresvanis/he-sample-operator/internal/conditions"
	"github.com/mresvanis/he-sample-operator/internal/finalizers"
	"github.com/mresvanis/he-sample-operator/internal/module"
	"github.com/mresvanis/he-sample-operator/internal/nodeselector"
)

const (
	testDeviceConfigName = "test"
)

var _ = Describe("DeviceConfigReconciler", func() {
	Describe("Reconcile", func() {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: testDeviceConfigName,
			},
		}

		Context("with a valid DeviceConfig", func() {
			ctx := context.TODO()
			dc := makeTestDeviceConfig()

			var (
				gCtrl *gomock.Controller
				mr    *module.MockReconciler
				fu    *finalizers.MockUpdater
				cu    *conditions.MockUpdater
				nsv   *nodeselector.MockValidator
				r     *DeviceConfigReconciler
				c     *client.MockClient
			)

			BeforeEach(func() {
				gCtrl = gomock.NewController(GinkgoT())
				mr = module.NewMockReconciler(gCtrl)
				fu = finalizers.NewMockUpdater(gCtrl)
				cu = conditions.NewMockUpdater(gCtrl)
				nsv = nodeselector.NewMockValidator(gCtrl)
				c = client.NewMockClient(gCtrl)
			})

			When("a client not-found error occurs", func() {
				BeforeEach(func() {
					s := scheme.Scheme

					r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, cu, nsv)

					gomock.InOrder(
						c.EXPECT().
							Get(ctx, req.NamespacedName, gomock.Any()).
							Return(apierrors.NewNotFound(schema.GroupResource{Resource: "deviceconfigs"}, testDeviceConfigName)),
					)
				})

				It("should not requeue or return an error", func() {
					res, err := r.Reconcile(ctx, req)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.Requeue).To(BeFalse())
				})
			})

			When("a client generic error occurs", func() {
				BeforeEach(func() {
					s := scheme.Scheme

					r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, cu, nsv)

					gomock.InOrder(
						c.EXPECT().
							Get(ctx, req.NamespacedName, gomock.Any()).
							Return(apierrors.NewServiceUnavailable("Service unavailable")),
					)
				})

				It("should not requeue and return an error", func() {
					res, err := r.Reconcile(ctx, req)
					Expect(err).To(HaveOccurred())
					Expect(res.Requeue).To(BeFalse())
				})
			})

			When("no client error occurs", func() {
				BeforeEach(func() {
					s := scheme.Scheme
					Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
					Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

					r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, cu, nsv)

					gomock.InOrder(
						c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
							func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
								d.ObjectMeta = dc.ObjectMeta
								d.Spec = dc.Spec
								return nil
							},
						),
						nsv.EXPECT().CheckDeviceConfigForConflictingNodeSelector(ctx, dc).Return(nil),
						fu.EXPECT().ContainsDeletionFinalizer(dc).Return(false),
						fu.EXPECT().AddDeletionFinalizer(ctx, dc).Return(nil),
						mr.EXPECT().ReconcileModule(ctx, dc).Return(nil),
						cu.EXPECT().SetConditionsReady(ctx, dc, "Reconciled", gomock.Any()).Return(nil),
					)
				})

				It("should create all resources", func() {
					res, err := r.Reconcile(ctx, req)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.Requeue).To(BeFalse())
				})
			})

			When("a reconcile Module error occurs", func() {
				BeforeEach(func() {
					s := scheme.Scheme
					Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
					Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

					r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, cu, nsv)

					gomock.InOrder(
						c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
							func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
								d.ObjectMeta = dc.ObjectMeta
								d.Spec = dc.Spec
								return nil
							},
						),
						nsv.EXPECT().CheckDeviceConfigForConflictingNodeSelector(ctx, dc).Return(nil),
						fu.EXPECT().ContainsDeletionFinalizer(dc).Return(false),
						fu.EXPECT().AddDeletionFinalizer(ctx, dc).Return(nil),
						mr.EXPECT().ReconcileModule(ctx, dc).Return(errors.New("some-error")),
						cu.EXPECT().SetConditionsErrored(ctx, dc, conditions.ReasonModuleFailed, gomock.Any()).Return(nil),
					)
				})

				It("should return the respective error", func() {
					res, err := r.Reconcile(ctx, req)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("some-error"))
					Expect(res.Requeue).To(BeFalse())
				})
			})

			Context("that does not contain a finalizer", func() {
				When("an add finalizer error occurs", func() {
					BeforeEach(func() {
						s := scheme.Scheme
						Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
						Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

						r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, cu, nsv)

						gomock.InOrder(
							c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
								func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
									d.ObjectMeta = dc.ObjectMeta
									d.Spec = dc.Spec
									return nil
								},
							),
							nsv.EXPECT().CheckDeviceConfigForConflictingNodeSelector(ctx, dc).Return(nil),
							fu.EXPECT().ContainsDeletionFinalizer(dc).Return(false),
							fu.EXPECT().AddDeletionFinalizer(ctx, dc).Return(errors.New("some-error")),
						)
					})

					It("should return the respective error", func() {
						res, err := r.Reconcile(ctx, req)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("some-error"))
						Expect(res.Requeue).To(BeFalse())
					})
				})
			})
		})

		Context("with a NodeSelectorConflictErrored DeviceConfig", func() {
			var (
				gCtrl        *gomock.Controller
				ctx          context.Context
				r            *DeviceConfigReconciler
				dc           *examplecomv1alpha1.DeviceConfig
				nsv          *nodeselector.MockValidator
				c            *client.MockClient
				fakeRecorder *record.FakeRecorder
			)

			BeforeEach(func() {
				gCtrl = gomock.NewController(GinkgoT())
				ctx = context.TODO()
				dc = makeTestDeviceConfig()
				nsv = nodeselector.NewMockValidator(gCtrl)
				c = client.NewMockClient(gCtrl)
			})

			It("should not return an error and record a conflicting selector event", func() {
				nsv.
					EXPECT().
					CheckDeviceConfigForConflictingNodeSelector(ctx, dc).
					Return(fmt.Errorf("an error"))

				s := scheme.Scheme
				Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())

				gomock.InOrder(
					c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
						func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
							d.ObjectMeta = dc.ObjectMeta
							d.Spec = dc.Spec
							return nil
						},
					),
				)

				fakeRecorder = record.NewFakeRecorder(1)
				r = NewDeviceConfigReconciler(c, s, fakeRecorder,
					module.NewReconciler(c, s),
					finalizers.NewUpdater(c),
					conditions.NewUpdater(c),
					nsv,
				)

				res, err := r.Reconcile(ctx, req)
				Expect(err).ToNot(HaveOccurred())
				Expect(res.Requeue).To(BeFalse())
				msg := <-fakeRecorder.Events
				Expect(msg).To(ContainSubstring("Conflicting DeviceConfig NodeSelectors found"))
			})
		})

		Context("with a deleted DeviceConfig", func() {
			ctx := context.TODO()
			dc := makeTestDeviceConfig(deletedAt(time.Now()))

			var (
				gCtrl *gomock.Controller
				mr    *module.MockReconciler
				fu    *finalizers.MockUpdater
				r     *DeviceConfigReconciler
				c     *client.MockClient
			)

			BeforeEach(func() {
				gCtrl = gomock.NewController(GinkgoT())
				mr = module.NewMockReconciler(gCtrl)
				fu = finalizers.NewMockUpdater(gCtrl)
				c = client.NewMockClient(gCtrl)
			})

			Context("which contains a deletion finalizer", func() {
				Context("and a deletion error occurs", func() {
					It("should return an error", func() {
						s := scheme.Scheme
						Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
						Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

						gomock.InOrder(
							c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
								func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
									d.ObjectMeta = dc.ObjectMeta
									d.Spec = dc.Spec
									return nil
								},
							),
						)

						r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, nil, nil)

						gomock.InOrder(
							fu.EXPECT().ContainsDeletionFinalizer(dc).Return(true),
							mr.EXPECT().DeleteModule(ctx, dc).Return(errors.New("something went wrong")),
						)

						res, err := r.Reconcile(ctx, req)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(And(
							ContainSubstring("failed to delete DeviceConfig resources"),
							ContainSubstring("something went wrong")))
						Expect(res.Requeue).To(BeFalse())
					})
				})

				Context("and no deletion error occurs", func() {
					Context("and no remove finalizer error occurs", func() {
						It("should not requeue or return an error", func() {
							s := scheme.Scheme
							Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
							Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

							gomock.InOrder(
								c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
									func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
										d.ObjectMeta = dc.ObjectMeta
										d.Spec = dc.Spec
										return nil
									},
								),
							)

							r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, nil, nil)

							gomock.InOrder(
								fu.EXPECT().ContainsDeletionFinalizer(dc).Return(true),
								mr.EXPECT().DeleteModule(ctx, dc).Return(nil),
								fu.EXPECT().RemoveDeletionFinalizer(ctx, dc).Return(nil),
							)

							res, err := r.Reconcile(ctx, req)
							Expect(err).ToNot(HaveOccurred())
							Expect(res.Requeue).To(BeFalse())
						})
					})

					Context("and a remove finalizer error occurs", func() {
						It("should not requeue and return an error", func() {
							s := scheme.Scheme
							Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
							Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

							gomock.InOrder(
								c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
									func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
										d.ObjectMeta = dc.ObjectMeta
										d.Spec = dc.Spec
										return nil
									},
								),
							)

							r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), mr, fu, nil, nil)

							gomock.InOrder(
								fu.EXPECT().ContainsDeletionFinalizer(dc).Return(true),
								mr.EXPECT().DeleteModule(ctx, dc).Return(nil),
								fu.EXPECT().RemoveDeletionFinalizer(ctx, dc).Return(errors.New("some error")),
							)

							res, err := r.Reconcile(ctx, req)
							Expect(err).To(HaveOccurred())
							Expect(res.Requeue).To(BeFalse())
						})
					})
				})
			})

			Context("which does not contain a deletion finalizer", func() {
				It("should do nothing", func() {
					gomock.InOrder(
						c.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).DoAndReturn(
							func(_ interface{}, _ interface{}, d *examplecomv1alpha1.DeviceConfig) error {
								d.ObjectMeta = dc.ObjectMeta
								d.Spec = dc.Spec
								return nil
							},
						),
						fu.EXPECT().ContainsDeletionFinalizer(dc).Return(false),
					)

					s := scheme.Scheme
					Expect(examplecomv1alpha1.AddToScheme(s)).ToNot(HaveOccurred())
					Expect(kmmv1beta1.AddToScheme(s)).ToNot(HaveOccurred())

					r = NewDeviceConfigReconciler(c, s, record.NewFakeRecorder(1), nil, fu, nil, nil)

					res, err := r.Reconcile(ctx, req)
					Expect(err).ToNot(HaveOccurred())
					Expect(res.Requeue).To(BeFalse())
				})
			})
		})
	})
})

func deletedAt(now time.Time) deviceConfigOptions {
	return func(c *examplecomv1alpha1.DeviceConfig) {
		wrapped := metav1.NewTime(now)
		c.ObjectMeta.DeletionTimestamp = &wrapped
	}
}

type deviceConfigOptions func(*examplecomv1alpha1.DeviceConfig)

func makeTestDeviceConfig(opts ...deviceConfigOptions) *examplecomv1alpha1.DeviceConfig {
	c := &examplecomv1alpha1.DeviceConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: testDeviceConfigName,
		},
		Spec: examplecomv1alpha1.DeviceConfigSpec{
			DriverImage:   "",
			DriverVersion: "",
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c
}
