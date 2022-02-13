//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LendingConfig) DeepCopyInto(out *LendingConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LendingConfig.
func (in *LendingConfig) DeepCopy() *LendingConfig {
	if in == nil {
		return nil
	}
	out := new(LendingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LendingConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LendingConfigList) DeepCopyInto(out *LendingConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LendingConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LendingConfigList.
func (in *LendingConfigList) DeepCopy() *LendingConfigList {
	if in == nil {
		return nil
	}
	out := new(LendingConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LendingConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LendingConfigSpec) DeepCopyInto(out *LendingConfigSpec) {
	*out = *in
	if in.Schedules != nil {
		in, out := &in.Schedules, &out.Schedules
		*out = make([]ScheduleConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LendingConfigSpec.
func (in *LendingConfigSpec) DeepCopy() *LendingConfigSpec {
	if in == nil {
		return nil
	}
	out := new(LendingConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LendingConfigStatus) DeepCopyInto(out *LendingConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LendingConfigStatus.
func (in *LendingConfigStatus) DeepCopy() *LendingConfigStatus {
	if in == nil {
		return nil
	}
	out := new(LendingConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleConfig) DeepCopyInto(out *ScheduleConfig) {
	*out = *in
	out.Start = in.Start
	out.End = in.End
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleConfig.
func (in *ScheduleConfig) DeepCopy() *ScheduleConfig {
	if in == nil {
		return nil
	}
	out := new(ScheduleConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleTiming) DeepCopyInto(out *ScheduleTiming) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleTiming.
func (in *ScheduleTiming) DeepCopy() *ScheduleTiming {
	if in == nil {
		return nil
	}
	out := new(ScheduleTiming)
	in.DeepCopyInto(out)
	return out
}
