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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChildprocessSpec defines the desired state of Childprocess
type ChildprocessSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Childprocess. Edit childprocess_types.go to remove/update
	//Foo string `json:"foo,omitempty"`

	Mpod v1.PodSpec `json:"mpod,omitempty"`

	Tpod string `json:"tpod,omitempty"`
}

// ChildprocessStatus defines the observed state of Childprocess
type ChildprocessStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Mpod string `json:"mpod,omitempty"`

	Tpod v1.PodStatus `json:"tpod,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Childprocess is the Schema for the childprocesses API
type Childprocess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChildprocessSpec   `json:"spec,omitempty"`
	Status ChildprocessStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChildprocessList contains a list of Childprocess
type ChildprocessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Childprocess `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Childprocess{}, &ChildprocessList{})
}
