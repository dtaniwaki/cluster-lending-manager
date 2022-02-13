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

type ScheduleTiming struct {
	// Time is ...
	Time string `json:"time,omitempty"`
	// DayOfWeek is ...
	DayOfWeek string `json:"dayOfWeek,omitempty"`
}

type ScheduleConfig struct {
	// StartAt is...
	Start ScheduleTiming `json:"start,omitempty"`
	// EndAt is...
	End ScheduleTiming `json:"end,omitempty"`
}

// LendingConfigSpec defines the desired state of LendingConfig
type LendingConfigSpec struct {
	Timezone string `json:"timezone,omitempty"`
	// Schedules is ...
	Schedules []ScheduleConfig `json:"schedules,omitempty"`
}

// LendingConfigStatus defines the observed state of LendingConfig
type LendingConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LendingConfig is the Schema for the LendingConfigs API
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
