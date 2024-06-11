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
	"os"
	"time"

	//nolint:golint
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ghrunnerv1 "ghrunner/selfhosted/api/v1"
)

var _ = Describe("GhRunner controller", func() {
	Context("GhRunner controller test", func() {

		const GhRunnerName = "test-ghrunner"

		ctx := context.Background()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      GhRunnerName,
				Namespace: GhRunnerName,
			},
		}

		typeNamespaceName := types.NamespacedName{
			Name:      GhRunnerName,
			Namespace: GhRunnerName,
		}
		ghrunner := &ghrunnerv1.GhRunner{}

		BeforeEach(func() {
			By("Creating the Namespace to perform the tests")
			err := k8sClient.Create(ctx, namespace)
			Expect(err).To(Not(HaveOccurred()))

			By("Setting the Image ENV VAR which stores the Operand image")
			err = os.Setenv("GHRUNNER_IMAGE", "example.com/image:test")
			Expect(err).To(Not(HaveOccurred()))

			By("creating the custom resource for the Kind GhRunner")
			err = k8sClient.Get(ctx, typeNamespaceName, ghrunner)
			if err != nil && errors.IsNotFound(err) {
				// Let's mock our custom resource at the same way that we would
				// apply on the cluster the manifest under config/samples
				ghrunner := &ghrunnerv1.GhRunner{
					ObjectMeta: metav1.ObjectMeta{
						Name:      GhRunnerName,
						Namespace: namespace.Name,
					},
					Spec: ghrunnerv1.GhRunnerSpec{
						Size: 1,
					},
				}

				err = k8sClient.Create(ctx, ghrunner)
				Expect(err).To(Not(HaveOccurred()))
			}
		})

		AfterEach(func() {
			By("removing the custom resource for the Kind GhRunner")
			found := &ghrunnerv1.GhRunner{}
			err := k8sClient.Get(ctx, typeNamespaceName, found)
			Expect(err).To(Not(HaveOccurred()))

			Eventually(func() error {
				return k8sClient.Delete(context.TODO(), found)
			}, 2*time.Minute, time.Second).Should(Succeed())

			// TODO(user): Attention if you improve this code by adding other context test you MUST
			// be aware of the current delete namespace limitations.
			// More info: https://book.kubebuilder.io/reference/envtest.html#testing-considerations
			By("Deleting the Namespace to perform the tests")
			_ = k8sClient.Delete(ctx, namespace)

			By("Removing the Image ENV VAR which stores the Operand image")
			_ = os.Unsetenv("GHRUNNER_IMAGE")
		})

		It("should successfully reconcile a custom resource for GhRunner", func() {
			By("Checking if the custom resource was successfully created")
			Eventually(func() error {
				found := &ghrunnerv1.GhRunner{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Reconciling the custom resource created")
			ghrunnerReconciler := &GhRunnerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := ghrunnerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if Deployment was successfully created in the reconciliation")
			Eventually(func() error {
				found := &appsv1.Deployment{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, time.Minute, time.Second).Should(Succeed())

			By("Checking the latest Status Condition added to the GhRunner instance")
			Eventually(func() error {
				if ghrunner.Status.Conditions != nil &&
					len(ghrunner.Status.Conditions) != 0 {
					latestStatusCondition := ghrunner.Status.Conditions[len(ghrunner.Status.Conditions)-1]
					expectedLatestStatusCondition := metav1.Condition{
						Type:   typeAvailableGhRunner,
						Status: metav1.ConditionTrue,
						Reason: "Reconciling",
						Message: fmt.Sprintf(
							"Deployment for custom resource (%s) with %d replicas created successfully",
							ghrunner.Name,
							ghrunner.Spec.Size),
					}
					if latestStatusCondition != expectedLatestStatusCondition {
						return fmt.Errorf("The latest status condition added to the GhRunner instance is not as expected")
					}
				}
				return nil
			}, time.Minute, time.Second).Should(Succeed())
		})
	})
})
