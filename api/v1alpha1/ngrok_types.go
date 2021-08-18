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
	"bytes"
	"html/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/utils"
)

// NgrokSpec defines the desired state of Ngrok
type NgrokSpec struct {
	Service string `json:"service"`
	Port    int32  `json:"port"`

	// +kubebuilder:validation:Enum=http;tcp
	// +kubebuilder:default:=http
	// +optional
	Protocol string `json:"protocol"`

	// +optional
	AuthToken string `json:"authtoken"`

	// +optional
	Auth string `json:"auth"`

	// +optional
	Hostname string `json:"hostname"`

	// +optional
	RemoteAddr string `json:"remote_addr"`

	// +kubebuilder:validation:Enum=us;eu;ap;au;sa;jp;in
	// +optional
	Region string `json:"region"`

	// +kubebuilder:default:=false
	// +optional
	Inspect bool `json:"inspect"`
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

func (n *Ngrok) GenerateConfiguration() (string, error) {
	var configuration bytes.Buffer

	templateEngine, err := template.New("ngrok").Parse(utils.TMPL)
	if err != nil {
		return "", err
	}

	err = templateEngine.Execute(&configuration, &n)
	if err != nil {
		return "", err
	}

	return configuration.String(), nil
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
