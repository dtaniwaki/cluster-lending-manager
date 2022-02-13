/*
Copyright 2022 Daisuke Taniwaki..

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LendingScheduleSpec defines the desired state of LendingSchedule
type LendingScheduleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LendingSchedule. Edit lendingschedule_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// LendingScheduleStatus defines the observed state of LendingSchedule
type LendingScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LendingSchedule is the Schema for the lendingschedules API
type LendingSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LendingScheduleSpec   `json:"spec,omitempty"`
	Status LendingScheduleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LendingScheduleList contains a list of LendingSchedule
type LendingScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LendingSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LendingSchedule{}, &LendingScheduleList{})
}
