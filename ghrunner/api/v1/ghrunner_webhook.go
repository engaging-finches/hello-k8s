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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var ghrunnerlog = logf.Log.WithName("ghrunner-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *GhRunner) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-ghrunner-ghrunner-v1-ghrunner,mutating=true,failurePolicy=fail,sideEffects=None,groups=ghrunner.ghrunner,resources=ghrunners,verbs=create;update,versions=v1,name=mghrunner.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &GhRunner{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *GhRunner) Default() {
	ghrunnerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-ghrunner-ghrunner-v1-ghrunner,mutating=false,failurePolicy=fail,sideEffects=None,groups=ghrunner.ghrunner,resources=ghrunners,verbs=create;update,versions=v1,name=vghrunner.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &GhRunner{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateCreate() (admission.Warnings, error) {
	ghrunnerlog.Info("validate create", "name", r.Name)
	ghrunnerlog.Info("validate create", "name", r.Spec.Repo)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, apierrors.NewBadRequest("Repo is required")
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	ghrunnerlog.Info("validate update", "name", r.Name)
	ghrunnerlog.Info("validate create", "name", r.Spec.Repo)

	// TODO(user): fill in your validation logic upon object update.
	return nil, apierrors.NewBadRequest("Repo is required")
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateDelete() (admission.Warnings, error) {
	ghrunnerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
