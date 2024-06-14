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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	//"strings"
)

// log is for logging in this package.
var ghrunnerlog = logf.Log.WithName("ghrunner-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *GhRunner) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

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

// getGitHubRegistrationToken makes a POST request to the GitHub API to get a registration token for the runner
func getGitHubRegistrationToken(accessToken, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runners/registration-token", owner, repo)

	// Create a new request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	// Set the necessary headers
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Parse the JSON response to extract the registration token
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("must provide valid owner, repo, and pat to get registration token. pat should have read/write administraction access to the repo")
	}
	return token, nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateCreate() (admission.Warnings, error) {
	name := r.Name
	owner := r.Spec.Owner
	repo := r.Spec.Repo
	pat := r.Spec.Pat
	ghrunnerlog.Info("validate create", "name", name)
	ghrunnerlog.Info("validate create", "owner", owner)
	ghrunnerlog.Info("validate create", "repo", repo)
	ghrunnerlog.Info("validate create", "pat", pat)

	_, err := getGitHubRegistrationToken(pat, owner, repo)
	if err != nil { //if token is invalid, then return an error
		return nil, apierrors.NewBadRequest(err.Error())
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	ghrunnerlog.Info("validate update", "name", r.Name)
	ghrunnerlog.Info("validate create", "name", r.Spec.Repo)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *GhRunner) ValidateDelete() (admission.Warnings, error) {
	ghrunnerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
