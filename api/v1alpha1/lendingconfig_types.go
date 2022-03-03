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

// Schedule is ...
type Schedule struct {
	// Start is ...
	Start *string `json:"start,omitempty"`
	// End is ...
	End *string `json:"end,omitempty"`
}

// DaySchedule is ...
type DaySchedule struct {
	// Hours is ...
	Hours []Schedule `json:"hours,omitempty"`
}

// Schedule is ...
type ScheduleSpec struct {
	// Default is ...
	Default *DaySchedule `json:"default,omitempty"`
	// Monday is ...
	Monday *DaySchedule `json:"monday,omitempty"`
	// Tuesday is ...
	Tuesday *DaySchedule `json:"tuesday,omitempty"`
	// Wednesday is ...
	Wednesday *DaySchedule `json:"wednesday,omitempty"`
	// Thursday is ...
	Thursday *DaySchedule `json:"thursday,omitempty"`
	// Friday is ...
	Friday *DaySchedule `json:"friday,omitempty"`
	// Saturday is ...
	Saturday *DaySchedule `json:"saturday,omitempty"`
	// Sunday is ...
	Sunday *DaySchedule `json:"sunday,omitempty"`
	// TODO: Support holidays.

	// Always is ...
	Always bool `json:"always,omitempty"`
}

// Target is ...
type Target struct {
	// APIVersion is ...
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is ...
	Kind string `json:"kind,omitempty"`
	// Name is ...
	Name *string `json:"name,omitempty"`
}

// LendingConfigSpec defines the desired state of LendingConfig
type LendingConfigSpec struct {
	// Timezone is ...
	Timezone string `json:"timezone,omitempty"`
	// Schedules is ...
	Schedule ScheduleSpec `json:"schedule,omitempty"`
	// TargetRefs is ...
	Targets []Target `json:"targets,omitempty"`
}

// ObjectReference contains enough information to let you locate the
// typed referenced object inside the same namespace.
type ObjectReference struct {
	// APIVersion is ...
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is ...
	Kind string `json:"kind,omitempty"`
	// Name is ...
	Name string `json:"name,omitempty"`
}

// LendingReference is ...
type LendingReference struct {
	// ObjectReference is ...
	ObjectReference `json:"objectReference,omitempty"`
	// Replicas is ...
	Replicas int64 `json:"replicas,omitempty"`
}

// LendingConfigStatus defines the observed state of LendingConfig
type LendingConfigStatus struct {
	// LendingReferences is ...
	LendingReferences []LendingReference `json:"objectReferences,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LendingConfig is the Schema for the lendingconfigs API
type LendingConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LendingConfigSpec   `json:"spec,omitempty"`
	Status LendingConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LendingConfigList contains a list of LendingConfig
type LendingConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LendingConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LendingConfig{}, &LendingConfigList{})
}
