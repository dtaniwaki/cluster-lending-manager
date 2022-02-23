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

	"github.com/dtaniwaki/cluster-lending-manager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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
	logger.Info(fmt.Sprintf("Start lending of %s/%s", cronctx.lendingconfig.Namespace, cronctx.lendingconfig.Name))

	for _, target := range cronctx.lendingconfig.Spec.Targets {
		groupVersionKind, err := getGroupVersionKind(target)
		if err != nil {
			return err
		}
		objs := &unstructured.Unstructured{}
		objs.SetGroupVersionKind(groupVersionKind)

		err = cronctx.reconciler.List(ctx, objs, &client.ListOptions{Namespace: cronctx.lendingconfig.Namespace})
		if err != nil {
			return err
		}
		err = objs.EachListItem(func(obj runtime.Object) error {
			uobj := obj.(*unstructured.Unstructured)
			logger.Info(fmt.Sprintf("Patch %s %s/%s", groupVersionKind, uobj.GetNamespace(), uobj.GetName()))

			replicas, found, err := unstructured.NestedInt64(uobj.UnstructuredContent(), "spec", "replicas")
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("The resource doesn't have replcias field.")
			}
			if replicas > 0 {
				logger.Info("Skipped the already running resource.")
				return nil
			}
			annotations := uobj.GetAnnotations()
			if annotations[annotationNameSkip] == "true" {
				logger.Info("Skipped the annotated resource.")
				return nil
			}
			// TODO: Restore the last replicas in the status.
			lastReplicas := 1
			patch := makeReplicasPatch(uobj, groupVersionKind, int32(lastReplicas))
			err = cronctx.reconciler.Patch(ctx, patch, client.Apply, &client.PatchOptions{
				FieldManager: "application/apply-patch",
				Force:        pointer.Bool(true),
			})
			if err != nil {
				return err
			}
			logger.Info("Patched the resource.")
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
	logger.Info(fmt.Sprintf("End lending of %s/%s", cronctx.lendingconfig.Namespace, cronctx.lendingconfig.Name))

	for _, target := range cronctx.lendingconfig.Spec.Targets {
		groupVersionKind, err := getGroupVersionKind(target)
		if err != nil {
			return err
		}
		objs := &unstructured.Unstructured{}
		objs.SetGroupVersionKind(groupVersionKind)

		err = cronctx.reconciler.List(ctx, objs, &client.ListOptions{Namespace: cronctx.lendingconfig.Namespace})
		if err != nil {
			return err
		}
		err = objs.EachListItem(func(obj runtime.Object) error {
			uobj := obj.(*unstructured.Unstructured)
			logger.Info(fmt.Sprintf("Patch %s %s/%s", groupVersionKind, uobj.GetNamespace(), uobj.GetName()))

			replicas, found, err := unstructured.NestedInt64(uobj.UnstructuredContent(), "spec", "replicas")
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("The resource doesn't have replcias field.")
			}
			if replicas == 0 {
				logger.Info("Skipped the already stopped resource.")
				return nil
			}
			annotations := uobj.GetAnnotations()
			if annotations[annotationNameSkip] == "true" {
				logger.Info("Skipped the annotated resource.")
				return nil
			}
			// TODO: Save the current replicas in the status.
			lastReplicas := 0
			patch := makeReplicasPatch(uobj, groupVersionKind, int32(lastReplicas))
			err = cronctx.reconciler.Patch(ctx, patch, client.Apply, &client.PatchOptions{
				FieldManager: "application/apply-patch",
				Force:        pointer.Bool(true),
			})
			if err != nil {
				return err
			}
			logger.Info("Patched the resource.")
			return nil
		})
		if err != nil {
			return err
		}
	}
	cronctx.reconciler.Recorder.Event(cronctx.lendingconfig.ToCompatible(), corev1.EventTypeNormal, LendingEnded, "Lending ended.")

	return nil
}

func getGroupVersionKind(obj v1alpha1.Target) (schema.GroupVersionKind, error) {
	groupVersion, err := schema.ParseGroupVersion(obj.APIVersion)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	return groupVersion.WithKind(obj.Kind), nil

}

func makeReplicasPatch(uobj *unstructured.Unstructured, groupVersionKind schema.GroupVersionKind, replicas int32) *unstructured.Unstructured {
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
