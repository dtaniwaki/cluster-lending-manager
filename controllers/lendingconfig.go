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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	clusterlendingmanagerv1alpha1 "github.com/dtaniwaki/cluster-lending-manager/api/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
)

type LendingConfig clusterlendingmanagerv1alpha1.LendingConfig

var hoursPattern = regexp.MustCompile(`(\d{2}):(\d{2}) *(am|pm)?`)

func (config *LendingConfig) ClearSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())
	return nil
}

func (config *LendingConfig) UpdateSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())

	items := []CronItem{}

	if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Default.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Monday != nil && config.Spec.Schedule.Monday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Monday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Tuesday != nil && config.Spec.Schedule.Tuesday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Tuesday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Wednesday != nil && config.Spec.Schedule.Wednesday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Wednesday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Thursday != nil && config.Spec.Schedule.Thursday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Thursday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Friday != nil && config.Spec.Schedule.Friday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Friday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Saturday != nil && config.Spec.Schedule.Saturday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Saturday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
		}
	}

	if config.Spec.Schedule.Sunday != nil && config.Spec.Schedule.Sunday.Hours != nil {
		crons, err := config.getCrons(reconciler, config.Spec.Schedule.Sunday.Hours)
		if err != nil {
			return err
		}
		for _, cron := range crons {
			items = append(items, cron)
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

func (config *LendingConfig) getCrons(reconciler *LendingConfigReconciler, schedules []clusterlendingmanagerv1alpha1.Schedule) ([]CronItem, error) {
	res := []CronItem{}
	for _, hours := range schedules {
		if hours.Start != nil {
			hour, minute, err := parseHours(*hours.Start)
			if err != nil {
				return nil, err
			}
			tsz := fmt.Sprintf("CRON_TZ=%s 0 %d %d * *", config.Spec.Timezone, minute, hour)
			res = append(res, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "LendingStart",
			)})
		}
		if hours.End != nil {
			hour, minute, err := parseHours(*hours.End)
			if err != nil {
				return nil, err
			}
			tsz := fmt.Sprintf("CRON_TZ=%s 0 %d %d * *", config.Spec.Timezone, minute, hour)
			res = append(res, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "LendingEnd",
			)})
		}
	}
	return res, nil
}

func parseHours(hours string) (int32, int32, error) {
	res := hoursPattern.FindStringSubmatch(strings.ToLower(hours))
	if len(res) != 3 {
		return 0, 0, errors.New("hours format is invalid")
	}
	hour, err := strconv.Atoi(res[0])
	if err != nil {
		return 0, 0, err
	}
	minute, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, 0, err
	}
	if res[2] == "pm" {
		hour += 12
	}
	return int32(hour), int32(minute), nil
}
