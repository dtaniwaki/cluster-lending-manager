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

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CronContext struct {
	reconciler    *LendingConfigReconciler
	lendingconfig *LendingConfig
	event         LendingEvent
}

func NewCronContext(
	reconciler *LendingConfigReconciler,
	lendingconfig *LendingConfig,
	event LendingEvent,
) *CronContext {
	return &CronContext{
		reconciler,
		lendingconfig,
		event,
	}
}

func (cronctx *CronContext) Run() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, CTX_VALUE_NAME, cronctx.lendingconfig.Name)
	ctx = context.WithValue(ctx, CTX_VALUE_NAMESPACE, cronctx.lendingconfig.Namespace)
	logger := log.FromContext(ctx)

	if err := cronctx.run(ctx); err != nil {
		logger.Error(err, "Failed to run a cron job")
	}
}

func (cronctx *CronContext) run(ctx context.Context) error {
	// TODO: Start or stop.
	return nil
}
