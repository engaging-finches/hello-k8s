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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GhRunnerSpec defines the desired state of GhRunner
type GhRunnerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size defines the number of GhRunner instances
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:ExclusiveMaximum=false
	Size  int32  `json:"size,omitempty"`
	Repo  string `json:"repo,omitempty"`
	Pat   string `json:"pat,omitempty"`
	Owner string `json:"owner,omitempty"`
}

// GhRunnerStatus defines the observed state of GhRunner
type GhRunnerStatus struct {
	// Represents the observations of a GhRunner's current state.
	// GhRunner.status.conditions.type are: "Available", "Progressing", and "Degraded"
	// GhRunner.status.conditions.status are one of True, False, Unknown.
	// GhRunner.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// GhRunner.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GhRunner is the Schema for the ghrunners API
type GhRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GhRunnerSpec   `json:"spec,omitempty"`
	Status GhRunnerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GhRunnerList contains a list of GhRunner
type GhRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GhRunner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GhRunner{}, &GhRunnerList{})
}
