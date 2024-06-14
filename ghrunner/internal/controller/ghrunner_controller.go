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
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ghrunnerv1 "ghrunner/selfhosted/api/v1"
)

const ghrunnerFinalizer = "ghrunner.ghrunner/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailableGhRunner represents the status of the Deployment reconciliation
	typeAvailableGhRunner = "Available"
	// typeDegradedGhRunner represents the status used when the custom resource is deleted and the finalizer operations are yet to occur.
	typeDegradedGhRunner = "Degraded"
)

// GhRunnerReconciler reconciles a GhRunner object
type GhRunnerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// The following markers are used to generate the rules permissions (RBAC) on config/rbac using controller-gen
// when the command <make manifests> is executed.
// To know more about markers see: https://book.kubebuilder.io/reference/markers.html

// +kubebuilder:rbac:groups=ghrunner.ghrunner,resources=ghrunners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ghrunner.ghrunner,resources=ghrunners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ghrunner.ghrunner,resources=ghrunners/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// It is essential for the controller's reconciliation loop to be idempotent. By following the Operator
// pattern you will create Controllers which provide a reconcile function
// responsible for synchronizing resources until the desired state is reached on the cluster.
// Breaking this recommendation goes against the design principles of controller-runtime.
// and may lead to unforeseen consequences such as resources becoming stuck and requiring manual intervention.
// For further info:
// - About Operator Pattern: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
// - About Controllers: https://kubernetes.io/docs/concepts/architecture/controller/
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile

func check_true(a bool, b bool, c bool, d bool) bool {
	if a || b || c || d {
		return true
	}
	return false
}

func (r *GhRunnerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the GhRunner instance
	// The purpose is check if the Custom Resource for the Kind GhRunner
	// is applied on the cluster if not we return nil to stop the reconciliation
	ghrunner := &ghrunnerv1.GhRunner{}
	err := r.Get(ctx, req.NamespacedName, ghrunner)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("ghrunner resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ghrunner")
		return ctrl.Result{}, err
	}

	// Log the new fields
	log.Info("Reconciling GhRunner", "Owner", ghrunner.Spec.Owner, "Repo", ghrunner.Spec.Repo, "PAT", "REDACTED")

	// Let's just set the status as Unknown when no status is available
	if ghrunner.Status.Conditions == nil || len(ghrunner.Status.Conditions) == 0 {
		meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeAvailableGhRunner, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
		if err = r.Status().Update(ctx, ghrunner); err != nil {
			log.Error(err, "Failed to update GhRunner status")
			return ctrl.Result{}, err
		}

		// Let's re-fetch the ghrunner Custom Resource after updating the status
		// so that we have the latest state of the resource on the cluster and we will avoid
		// raising the error "the object has been modified, please apply
		// your changes to the latest version and try again" which would re-trigger the reconciliation
		// if we try to update it again in the following operations
		if err := r.Get(ctx, req.NamespacedName, ghrunner); err != nil {
			log.Error(err, "Failed to re-fetch ghrunner")
			return ctrl.Result{}, err
		}
	}

	// Let's add a finalizer. Then, we can define some operations which should
	// occur before the custom resource is deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(ghrunner, ghrunnerFinalizer) {
		log.Info("Adding Finalizer for GhRunner")
		if ok := controllerutil.AddFinalizer(ghrunner, ghrunnerFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, ghrunner); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Check if the GhRunner instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isGhRunnerMarkedToBeDeleted := ghrunner.GetDeletionTimestamp() != nil
	if isGhRunnerMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(ghrunner, ghrunnerFinalizer) {
			log.Info("Performing Finalizer Operations for GhRunner before delete CR")

			// Let's add here a status "Downgrade" to reflect that this resource began its process to be terminated.
			meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeDegradedGhRunner,
				Status: metav1.ConditionUnknown, Reason: "Finalizing",
				Message: fmt.Sprintf("Performing finalizer operations for the custom resource: %s ", ghrunner.Name)})

			if err := r.Status().Update(ctx, ghrunner); err != nil {
				log.Error(err, "Failed to update GhRunner status")
				return ctrl.Result{}, err
			}

			// Perform all operations required before removing the finalizer and allow
			// the Kubernetes API to remove the custom resource.
			r.doFinalizerOperationsForGhRunner(ghrunner)

			// TODO(user): If you add operations to the doFinalizerOperationsForGhRunner method
			// then you need to ensure that all worked fine before deleting and updating the Downgrade status
			// otherwise, you should requeue here.

			// Re-fetch the ghrunner Custom Resource before updating the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raising the error "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, ghrunner); err != nil {
				log.Error(err, "Failed to re-fetch ghrunner")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeDegradedGhRunner,
				Status: metav1.ConditionTrue, Reason: "Finalizing",
				Message: fmt.Sprintf("Finalizer operations for custom resource %s name were successfully accomplished", ghrunner.Name)})

			if err := r.Status().Update(ctx, ghrunner); err != nil {
				log.Error(err, "Failed to update GhRunner status")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for GhRunner after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(ghrunner, ghrunnerFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for GhRunner")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, ghrunner); err != nil {
				log.Error(err, "Failed to remove finalizer for GhRunner")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: ghrunner.Name, Namespace: ghrunner.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForGhRunner(ghrunner)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for GhRunner")

			// The following implementation will update the status
			meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeAvailableGhRunner,
				Status: metav1.ConditionFalse, Reason: "Reconciling",
				Message: fmt.Sprintf("Failed to create Deployment for the custom resource (%s): (%s)", ghrunner.Name, err)})

			if err := r.Status().Update(ctx, ghrunner); err != nil {
				log.Error(err, "Failed to update GhRunner status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}

	// The CRD API defines that the GhRunner type have a GhRunnerSpec.Size field
	// to set the quantity of Deployment instances to the desired state on the cluster.
	// Therefore, the following code will ensure the Deployment size is the same as defined
	// via the Size spec of the Custom Resource which we are reconciling.
	size := ghrunner.Spec.Size
	repo := ghrunner.Spec.Repo
	owner := ghrunner.Spec.Owner
	pat := ghrunner.Spec.Pat
	log.Info("Checking the size of the Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name, "Size", size)
	log.Info("Checking the owner of the Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name, "Owner", owner)
	log.Info("Checking the repo of the Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name, "Repo", repo)
	// log.Info("Checking the PAT of the Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name, "PAT", "REDACTED")

	size_check := *found.Spec.Replicas != size
	owner_check := found.Spec.Template.Spec.Containers[0].Env[0].Value != owner
	repo_check := found.Spec.Template.Spec.Containers[0].Env[1].Value != repo
	pat_check := found.Spec.Template.Spec.Containers[0].Env[2].Value != pat

	if check_true(size_check, owner_check, repo_check, pat_check) {
		found.Spec.Replicas = &size
		found.Spec.Template.Spec.Containers[0].Env[0].Value = owner
		found.Spec.Template.Spec.Containers[0].Env[1].Value = repo
		found.Spec.Template.Spec.Containers[0].Env[2].Value = pat
		if err = r.Update(ctx, found); err != nil {
			log.Error(err, "Failed to update Deployment",
				"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

			// Re-fetch the ghrunner Custom Resource before updating the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raising the error "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, ghrunner); err != nil {
				log.Error(err, "Failed to re-fetch ghrunner")
				return ctrl.Result{}, err
			}

			// The following implementation will update the status
			meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeAvailableGhRunner,
				Status: metav1.ConditionFalse, Reason: "Resizing",
				Message: fmt.Sprintf("Failed to update the size for the custom resource (%s): (%s)", ghrunner.Name, err)})

			if err := r.Status().Update(ctx, ghrunner); err != nil {
				log.Error(err, "Failed to update GhRunner status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		// Now, that we update the size we want to requeue the reconciliation
		// so that we can ensure that we have the latest state of the resource before
		// update. Also, it will help ensure the desired state on the cluster
		return ctrl.Result{Requeue: true}, nil
	}

	// The following implementation will update the status
	meta.SetStatusCondition(&ghrunner.Status.Conditions, metav1.Condition{Type: typeAvailableGhRunner,
		Status: metav1.ConditionTrue, Reason: "Reconciling",
		Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", ghrunner.Name, size)})

	if err := r.Status().Update(ctx, ghrunner); err != nil {
		log.Error(err, "Failed to update GhRunner status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// finalizeGhRunner will perform the required operations before delete the CR.
func (r *GhRunnerReconciler) doFinalizerOperationsForGhRunner(cr *ghrunnerv1.GhRunner) {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.

	// Note: It is not recommended to use finalizers with the purpose of deleting resources which are
	// created and managed in the reconciliation. These ones, such as the Deployment created on this reconcile,
	// are defined as dependent of the custom resource. See that we use the method ctrl.SetControllerReference.
	// to set the ownerRef which means that the Deployment will be deleted by the Kubernetes API.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/use-cascading-deletion/

	// The following implementation will raise an event
	r.Recorder.Event(cr, "Warning", "Deleting",
		fmt.Sprintf("Custom Resource %s is being deleted from the namespace %s",
			cr.Name,
			cr.Namespace))
}

// deploymentForGhRunner returns a GhRunner Deployment object
func (r *GhRunnerReconciler) deploymentForGhRunner(
	ghrunner *ghrunnerv1.GhRunner) (*appsv1.Deployment, error) {
	ls := labelsForGhRunner(ghrunner.Name)
	replicas := ghrunner.Spec.Size
	pat := ghrunner.Spec.Pat
	owner := ghrunner.Spec.Owner
	repo := ghrunner.Spec.Repo

	// Get the Operand image
	image, err := imageForGhRunner()
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ghrunner.Name,
			Namespace: ghrunner.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					// TODO(user): Uncomment the following code to configure the nodeAffinity expression
					// according to the platforms which are supported by your solution. It is considered
					// best practice to support multiple architectures. build your manager image using the
					// makefile target docker-buildx. Also, you can use docker manifest inspect <image>
					// to check what are the platforms supported.
					// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
					//Affinity: &corev1.Affinity{
					//	NodeAffinity: &corev1.NodeAffinity{
					//		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					//			NodeSelectorTerms: []corev1.NodeSelectorTerm{
					//				{
					//					MatchExpressions: []corev1.NodeSelectorRequirement{
					//						{
					//							Key:      "kubernetes.io/arch",
					//							Operator: "In",
					//							Values:   []string{"amd64", "arm64", "ppc64le", "s390x"},
					//						},
					//						{
					//							Key:      "kubernetes.io/os",
					//							Operator: "In",
					//							Values:   []string{"linux"},
					//						},
					//					},
					//				},
					//			},
					//		},
					//	},
					//},
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &[]bool{true}[0],
						// IMPORTANT: seccomProfile was introduced with Kubernetes 1.19
						// If you are looking for to produce solutions to be supported
						// on lower versions you must remove this option.
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "ghrunner",
						ImagePullPolicy: corev1.PullIfNotPresent,
						// Ensure restrictive context for the container
						// More info: https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot:             &[]bool{true}[0],
							RunAsUser:                &[]int64{1000}[0],
							AllowPrivilegeEscalation: &[]bool{false}[0],
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{
									"ALL",
								},
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "REPO_OWNER",
								Value: owner,
							},
							{
								Name:  "REPO_NAME",
								Value: repo,
							},
							{
								Name:  "PAT",
								Value: pat,
							},
						},
					}},
				},
			},
		},
	}

	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(ghrunner, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// labelsForGhRunner returns the labels for selecting the resources
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func labelsForGhRunner(name string) map[string]string {
	fmt.Print(name)
	var imageTag string
	image, err := imageForGhRunner()
	if err == nil {
		imageTag = strings.Split(image, ":")[1]
	}
	return map[string]string{"app.kubernetes.io/name": "ghrunner",
		"app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/managed-by": "GhRunnerController",
	}
}

// imageForGhRunner gets the Operand image which is managed by this controller
// from the GHRUNNER_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForGhRunner() (string, error) {
	var imageEnvVar = "GHRUNNER_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

// SetupWithManager sets up the controller with the Manager.
// Note that the Deployment will be also watched in order to ensure its
// desirable state on the cluster
func (r *GhRunnerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ghrunnerv1.GhRunner{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
