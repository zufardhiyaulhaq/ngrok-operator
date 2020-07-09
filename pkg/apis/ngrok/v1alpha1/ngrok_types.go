package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NgrokSpec defines the desired state of Ngrok
type NgrokSpec struct {
	Service string `json:"service"`
	Port    int32  `json:"port"`
}

// NgrokStatus defines the observed state of Ngrok
type NgrokStatus struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ngrok is the Schema for the ngroks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ngroks,scope=Namespaced
type Ngrok struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NgrokSpec   `json:"spec,omitempty"`
	Status NgrokStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NgrokList contains a list of Ngrok
type NgrokList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ngrok `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ngrok{}, &NgrokList{})
}
