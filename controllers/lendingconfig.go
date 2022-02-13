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

package controllers

import (
	"context"
	clusterlendingmanagerv1alpha1 "dtaniwaki/cluster-lending-manager/api/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
)

type LendingConfig clusterlendingmanagerv1alpha1.LendingConfig

func (config *LendingConfig) ClearSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	return nil
}

func (config *LendingConfig) UpdateSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	return nil
}

func (config *LendingConfig) ToCompatible() *clusterlendingmanagerv1alpha1.LendingConfig {
	return (*clusterlendingmanagerv1alpha1.LendingConfig)(config)
}

func (config *LendingConfig) ToNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: config.Namespace, Name: config.Name}
}
