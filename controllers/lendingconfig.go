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
	"fmt"

	clusterlendingmanagerv1alpha1 "github.com/dtaniwaki/cluster-lending-manager/api/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
)

type LendingConfig clusterlendingmanagerv1alpha1.LendingConfig

func (config *LendingConfig) ClearSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())
	return nil
}

func (config *LendingConfig) UpdateSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())

	items := []CronItem{}

	if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Default.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Monday != nil && config.Spec.Schedule.Monday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Monday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Tuesday != nil && config.Spec.Schedule.Tuesday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Tuesday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Wednesday != nil && config.Spec.Schedule.Wednesday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Wednesday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Thursday != nil && config.Spec.Schedule.Thursday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Thursday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Friday != nil && config.Spec.Schedule.Friday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Friday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Saturday != nil && config.Spec.Schedule.Saturday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Saturday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	if config.Spec.Schedule.Sunday != nil && config.Spec.Schedule.Sunday.Hours != nil {
		for _, tsz := range getCrons(config.Spec.Timezone, config.Spec.Schedule.Sunday.Hours) {
			items = append(items, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "",
			)})
		}
	}

	err := reconciler.Cron.Add(config.ToNamespacedName(), items)
	if err != nil {
		return err
	}

	return nil
}

func (config *LendingConfig) ToCompatible() *clusterlendingmanagerv1alpha1.LendingConfig {
	return (*clusterlendingmanagerv1alpha1.LendingConfig)(config)
}

func (config *LendingConfig) ToNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: config.Namespace, Name: config.Name}
}

func parseHours(hours string) (int32, int32) {
	return 0, 0
}

func getCrons(timezone string, schedules []clusterlendingmanagerv1alpha1.Schedule) []string {
	res := []string{}
	for _, hours := range schedules {
		if hours.Start != nil {
			hour, minute := parseHours(*hours.Start)
			tsz := fmt.Sprintf("CRON_TZ=%s 0 %d %d * *", timezone, minute, hour)
			res = append(res, tsz)
		}
		if hours.End != nil {
			hour, minute := parseHours(*hours.End)
			tsz := fmt.Sprintf("CRON_TZ=%s 0 %d %d * *", timezone, minute, hour)
			res = append(res, tsz)
		}
	}
	return res
}
