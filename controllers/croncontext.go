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
	"sync"

	"github.com/dtaniwaki/cluster-lending-manager/api/v1alpha1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type GlobalCronContext struct {
	mu sync.Mutex
}

var gCronContext *GlobalCronContext

func init() {
	gCronContext = &GlobalCronContext{}
}

type CronContext struct {
	reconciler    *LendingConfigReconciler
	lendingConfig *LendingConfig
	event         LendingEvent
}

func NewCronContext(
	reconciler *LendingConfigReconciler,
	lendingConfig *LendingConfig,
	event LendingEvent,
) *CronContext {
	return &CronContext{
		reconciler,
		lendingConfig,
		event,
	}
}

func (cronctx *CronContext) Run() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, CTX_VALUE_NAME, cronctx.lendingConfig.Name)
	ctx = context.WithValue(ctx, CTX_VALUE_NAMESPACE, cronctx.lendingConfig.Namespace)
	logger := log.FromContext(ctx)

	if err := cronctx.run(ctx); err != nil {
		logger.Error(err, "Failed to run a cron job")
	}
}

func (cronctx *CronContext) run(ctx context.Context) error {
	gCronContext.mu.Lock()
	defer gCronContext.mu.Unlock()

	if cronctx.event == LendingStart {
		return cronctx.startLending(ctx)
	} else if cronctx.event == LendingEnd {
		return cronctx.endLending(ctx)
	} else {
		return errors.New(fmt.Sprintf("Unknown event %s", cronctx.event))
	}
}

func (cronctx *CronContext) startLending(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("Start lending")

	lendingConfig := &LendingConfig{}
	err := cronctx.reconciler.Get(ctx, cronctx.lendingConfig.ToNamespacedName(), lendingConfig.ToCompatible())
	if err != nil {
		return errors.WithStack(err)
	}

	err = lendingConfig.ActivateTargetResources(ctx, cronctx.reconciler)
	if err != nil {
		return err
	}

	lendingConfig.Status.LendingReferences = []v1alpha1.LendingReference{}

	err = cronctx.reconciler.Status().Update(ctx, lendingConfig.ToCompatible(), &client.UpdateOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	cronctx.reconciler.Recorder.Event(lendingConfig.ToCompatible(), corev1.EventTypeNormal, LendingStarted, "Lending started.")

	return nil
}

func (cronctx *CronContext) endLending(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("End lending")

	lendingConfig := &LendingConfig{}
	err := cronctx.reconciler.Get(ctx, cronctx.lendingConfig.ToNamespacedName(), lendingConfig.ToCompatible())
	if err != nil {
		return errors.WithStack(err)
	}

	lendingReferences, err := lendingConfig.DeactivateTargetResources(ctx, cronctx.reconciler)
	if err != nil {
		return err
	}

	lendingConfig.Status.LendingReferences = lendingReferences
	if err != nil {
		return err
	}

	err = cronctx.reconciler.Status().Update(ctx, lendingConfig.ToCompatible(), &client.UpdateOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	cronctx.reconciler.Recorder.Event(lendingConfig.ToCompatible(), corev1.EventTypeNormal, LendingEnded, "Lending ended.")

	return nil
}

func getGroupVersionKind(obj v1alpha1.Target) (schema.GroupVersionKind, error) {
	groupVersion, err := schema.ParseGroupVersion(obj.APIVersion)
	if err != nil {
		return schema.GroupVersionKind{}, errors.WithStack(err)
	}
	return groupVersion.WithKind(obj.Kind), nil

}

func makeReplicasPatch(uobj *unstructured.Unstructured, groupVersionKind schema.GroupVersionKind, replicas int64) *unstructured.Unstructured {
	patch := &unstructured.Unstructured{}
	// NOTE: obj.GetObjectKind().GroupVersionKind() is empty unexpextedly.
	patch.SetGroupVersionKind(groupVersionKind)
	patch.SetNamespace(uobj.GetNamespace())
	patch.SetName(uobj.GetName())
	patch.UnstructuredContent()["spec"] = map[string]interface{}{
		"replicas": replicas,
	}
	return patch
}
