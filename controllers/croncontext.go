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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	if cronctx.event == LendingStart {
		return cronctx.startLending(ctx)
	} else if cronctx.event == LendingEnd {
		return cronctx.endLending(ctx)
	} else {
		return fmt.Errorf("Unknown event %s", cronctx.event)
	}
}

func (cronctx *CronContext) startLending(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("Start lending")

	for _, target := range cronctx.lendingconfig.Spec.Targets {
		groupVersion, err := schema.ParseGroupVersion(target.APIVersion)
		if err != nil {
			return err
		}
		objs := &unstructured.Unstructured{}
		objs.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   groupVersion.Group,
			Version: groupVersion.Version,
			Kind:    target.Kind,
		})

		err = cronctx.reconciler.List(ctx, objs, &client.ListOptions{Namespace: cronctx.lendingconfig.Namespace})
		if err != nil {
			return err
		}
		err = objs.EachListItem(func(obj runtime.Object) error {
			metaobj := obj.(metav1.Object)
			patch := &unstructured.Unstructured{}
			patch.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
			patch.SetNamespace(metaobj.GetNamespace())
			patch.SetName(metaobj.GetName())
			logger.Info(fmt.Sprintf("Patch %s %s/%s", obj.GetObjectKind().GroupVersionKind(), metaobj.GetNamespace(), metaobj.GetName()))
			patch.UnstructuredContent()["spec"] = map[string]interface{}{
				"replicas": pointer.Int32(1),
			}
			err := cronctx.reconciler.Patch(ctx, patch, client.Apply, &client.PatchOptions{
				Force: pointer.Bool(true),
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	cronctx.reconciler.Recorder.Event(cronctx.lendingconfig.ToCompatible(), corev1.EventTypeNormal, LendingStarted, "Lending started.")

	return nil
}

func (cronctx *CronContext) endLending(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("End lending")

	for _, target := range cronctx.lendingconfig.Spec.Targets {
		groupVersion, err := schema.ParseGroupVersion(target.APIVersion)
		if err != nil {
			return err
		}
		patch := &unstructured.Unstructured{}
		patch.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   groupVersion.Group,
			Version: groupVersion.Version,
			Kind:    target.Kind,
		})
		patch.SetNamespace(cronctx.lendingconfig.Namespace)
		if target.Name != nil {
			patch.SetName(*target.Name)
		}
		patch.UnstructuredContent()["spec"] = map[string]interface{}{
			"replicas": pointer.Int32(0),
		}

		logger.Info(fmt.Sprintf("Patch %s %s/%s", patch.GetObjectKind().GroupVersionKind(), patch.GetNamespace(), patch.GetName()))
		err = cronctx.reconciler.Patch(ctx, patch, client.Apply, &client.PatchOptions{
			Force: pointer.Bool(true),
		})
		if err != nil {
			return err
		}
	}

	cronctx.reconciler.Recorder.Event(cronctx.lendingconfig.ToCompatible(), corev1.EventTypeNormal, LendingEnded, "Lending ended.")

	return nil
}
