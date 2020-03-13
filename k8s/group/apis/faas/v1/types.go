package v1

import (
	api_core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Function todo
type Function struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec FunctionSpec `json:"spec"`
	// +optional
	Status FunctionStatus `json:"status"`
}

// FunctionSpec todo
type FunctionSpec struct {
	Image     string                            `json:"image"`
	Replicas  int                               `json:"replicas"`
	Resources *api_core_v1.ResourceRequirements `json:"resources"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FunctionList todo
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Function `json:"items"`
}

// FunctionStatus todo
type FunctionStatus struct {
	Blah string
}
