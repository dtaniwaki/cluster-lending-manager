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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type LendingConfig clusterlendingmanagerv1alpha1.LendingConfig

var hoursPattern = regexp.MustCompile(`(\d{2}):(\d{2}) *(am|pm)?`)

type LendingConfigEvent = string

const annotationNameSkip = "cron-hpa.dtaniwaki.github.com/skip"

const (
	SchedulesUpdated LendingConfigEvent = "SchedulesUpdated"
	SchedulesCleared LendingConfigEvent = "SchedulesCleared"
	LendingStarted   LendingConfigEvent = "LendingStarted"
	LendingEnded     LendingConfigEvent = "endingEnded"
)

func (config *LendingConfig) ClearSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())

	reconciler.Recorder.Event(config.ToCompatible(), corev1.EventTypeNormal, SchedulesCleared, "Schedules cleared.")

	return nil
}

func (config *LendingConfig) UpdateSchedules(ctx context.Context, reconciler *LendingConfigReconciler) error {
	reconciler.Cron.Clear(config.ToNamespacedName())

	items := []CronItem{}

	var monHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Monday != nil && config.Spec.Schedule.Monday.Hours != nil {
		monHours = config.Spec.Schedule.Monday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		monHours = config.Spec.Schedule.Default.Hours
	}
	if monHours != nil {
		crons, err := config.getCrons(reconciler, "mon", monHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var tueHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Tuesday != nil && config.Spec.Schedule.Tuesday.Hours != nil {
		tueHours = config.Spec.Schedule.Tuesday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		tueHours = config.Spec.Schedule.Default.Hours
	}
	if tueHours != nil {
		crons, err := config.getCrons(reconciler, "tue", tueHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var wedHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Wednesday != nil && config.Spec.Schedule.Wednesday.Hours != nil {
		wedHours = config.Spec.Schedule.Wednesday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		wedHours = config.Spec.Schedule.Default.Hours
	}
	if wedHours != nil {
		crons, err := config.getCrons(reconciler, "wed", wedHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var thuHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Thursday != nil && config.Spec.Schedule.Thursday.Hours != nil {
		thuHours = config.Spec.Schedule.Thursday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		thuHours = config.Spec.Schedule.Default.Hours
	}
	if thuHours != nil {
		crons, err := config.getCrons(reconciler, "thu", thuHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var friHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Friday != nil && config.Spec.Schedule.Friday.Hours != nil {
		friHours = config.Spec.Schedule.Friday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		friHours = config.Spec.Schedule.Default.Hours
	}
	if friHours != nil {
		crons, err := config.getCrons(reconciler, "fri", friHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var satHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Saturday != nil && config.Spec.Schedule.Saturday.Hours != nil {
		satHours = config.Spec.Schedule.Saturday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		satHours = config.Spec.Schedule.Default.Hours
	}
	if satHours != nil {
		crons, err := config.getCrons(reconciler, "sat", satHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	var sunHours []clusterlendingmanagerv1alpha1.Schedule
	if config.Spec.Schedule.Sunday != nil && config.Spec.Schedule.Sunday.Hours != nil {
		sunHours = config.Spec.Schedule.Sunday.Hours
	} else if config.Spec.Schedule.Default != nil && config.Spec.Schedule.Default.Hours != nil {
		sunHours = config.Spec.Schedule.Default.Hours
	}
	if sunHours != nil {
		crons, err := config.getCrons(reconciler, "sun", sunHours)
		if err != nil {
			return err
		}
		items = append(items, crons...)
	}

	err := reconciler.Cron.Add(config.ToNamespacedName(), items)
	if err != nil {
		return err
	}

	reconciler.Recorder.Event(config.ToCompatible(), corev1.EventTypeNormal, SchedulesUpdated, "Schedules updated.")

	return nil
}

func (config *LendingConfig) ToCompatible() *clusterlendingmanagerv1alpha1.LendingConfig {
	return (*clusterlendingmanagerv1alpha1.LendingConfig)(config)
}

func (config *LendingConfig) ToNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: config.Namespace, Name: config.Name}
}

func (config *LendingConfig) getCrons(reconciler *LendingConfigReconciler, dayOfWeek string, schedules []clusterlendingmanagerv1alpha1.Schedule) ([]CronItem, error) {
	res := []CronItem{}
	for _, hours := range schedules {
		if hours.Start != nil {
			hour, minute, err := parseHours(*hours.Start)
			if err != nil {
				return nil, err
			}
			tsz := fmt.Sprintf("CRON_TZ=%s %d %d * * %s", config.Spec.Timezone, minute, hour, dayOfWeek)
			res = append(res, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "LendingStart",
			)})
		}
		if hours.End != nil {
			hour, minute, err := parseHours(*hours.End)
			if err != nil {
				return nil, err
			}
			tsz := fmt.Sprintf("CRON_TZ=%s %d %d * * %s", config.Spec.Timezone, minute, hour, dayOfWeek)
			res = append(res, CronItem{Cron: tsz, Job: NewCronContext(
				reconciler, config, "LendingEnd",
			)})
		}
	}
	return res, nil
}

func parseHours(hours string) (int32, int32, error) {
	res := hoursPattern.FindStringSubmatch(strings.ToLower(hours))
	if len(res) != 4 {
		return 0, 0, errors.New("hours format is invalid")
	}
	hour, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, 0, err
	}
	minute, err := strconv.Atoi(res[2])
	if err != nil {
		return 0, 0, err
	}
	if res[3] == "pm" {
		hour += 12
	}
	return int32(hour), int32(minute), nil
}
