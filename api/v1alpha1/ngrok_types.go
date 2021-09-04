/*
Copyright 2021.

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

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NgrokSpec defines the desired state of Ngrok
type NgrokSpec struct {
	Service string `json:"service"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// +kubebuilder:validation:Enum=http;tcp;tls
	// +kubebuilder:default:=http
	// +optional
	Protocol string `json:"protocol"`

	// +optional
	AuthToken string `json:"authtoken"`

	// +kubebuilder:validation:Enum=plain;secret
	// +kubebuilder:default:=plain
	// +optional
	AuthTokenType string `json:"authtoken_type"`

	// +kubebuilder:validation:Enum=us;eu;ap;au;sa;jp;in
	// +optional
	Region string `json:"region"`

	// +optional
	Auth string `json:"auth"`

	// +optional
	HostHeader string `json:"host_header"`

	// +kubebuilder:validation:Enum=true;false;both
	// +optional
	BindTLS string `json:"bind_tls"`

	// +kubebuilder:validation:Enum=true;false
	// +kubebuilder:default:=false
	// +optional
	Inspect bool `json:"inspect"`

	// +optional
	Hostname string `json:"hostname"`

	// +optional
	RemoteAddr string `json:"remote_addr"`

	// +kubebuilder:default:={image: zufardhiyaulhaq/ngrok}
	// +optional
	PodSpec PodSpec `json:"podSpec"`
}

func (n *NgrokSpec) Validate() error {
	if n.Protocol == "http" || n.Protocol == "tls" {
		if n.Service == "" {
			return fmt.Errorf("service invalid")
		}
	}

	if n.AuthToken == "" && n.Protocol == "tls" {
		return fmt.Errorf("protocol TLS only available in pro")
	}

	if n.Protocol == "tcp" {
		if n.RemoteAddr == "" {
			return fmt.Errorf("remote_addr invalid")
		}
	}

	return nil
}

type PodSpec struct {
	// +optional
	Image string `json:"image"`
}

// NgrokStatus defines the observed state of Ngrok
type NgrokStatus struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// Ngrok is the Schema for the Ngrok API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="Ngrok status"
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url",description="Ngrok URL"
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
