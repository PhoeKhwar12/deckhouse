/*
Copyright 2021 Flant JSC

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

package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"
)

const extendedMonitoringAnnotationKey = "extended-monitoring.flant.com/enabled"

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:                         "namespaces",
			ApiVersion:                   "v1",
			Kind:                         "Namespace",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "deployments",
			ApiVersion:                   "apps/v1",
			Kind:                         "Deployment",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "statefulsets",
			ApiVersion:                   "apps/v1",
			Kind:                         "StatefulSet",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "daemonsets",
			ApiVersion:                   "apps/v1",
			Kind:                         "DaemonSet",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "cronjobs",
			ApiVersion:                   "batch/v1beta1",
			Kind:                         "CronJob",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "ingresses",
			ApiVersion:                   "networking.k8s.io/v1",
			Kind:                         "Ingress",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:                         "nodes",
			ApiVersion:                   "v1",
			Kind:                         "Node",
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyNameNamespaceFilter,
		},
		{
			Name:       "migrated",
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			NameSelector: &types.NameSelector{
				MatchNames: []string{"d8-extended-monitoring-migrated"},
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-monitoring"},
				},
			},
			ExecuteHookOnSynchronization: pointer.Bool(false),
			ExecuteHookOnEvents:          pointer.Bool(false),
			FilterFunc:                   applyMigratedFilter,
		},
	},
}, handleLegacyAnnotatedIngress)

type ObjectNameNamespace struct {
	Name      string
	Namespace string
}

func applyNameNamespaceFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	if _, ok := obj.GetAnnotations()[extendedMonitoringAnnotationKey]; !ok {
		return nil, nil
	}

	return &ObjectNameNamespace{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}, nil
}

func applyMigratedFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return obj.GetName(), nil
}

func handleLegacyAnnotatedIngress(input *go_hook.HookInput) error {
	if len(input.Snapshots["migrated"]) > 0 {
		return nil
	}

	for _, obj := range input.Snapshots["namespaces"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "v1", "Namespace",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["deployments"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "apps/v1", "Deployment",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["statefulsets"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "apps/v1", "StatefulSet",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["daemonsets"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "apps/v1", "DaemonSet",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["cronjobs"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "batch/v1beta1", "CronJob",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["ingresses"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "networking.k8s.io/v1", "Ingress",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	for _, obj := range input.Snapshots["nodes"] {
		if obj == nil {
			continue
		}

		objMeta := obj.(*ObjectNameNamespace)

		input.PatchCollector.Filter(filterLabelsAndAnnotations, "v1", "Node",
			objMeta.Namespace, objMeta.Name, object_patch.IgnoreMissingObject())
	}

	return nil
}

func filterLabelsAndAnnotations(originalObj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	obj := originalObj.DeepCopy()

	annotations := obj.GetAnnotations()
	labels := obj.GetLabels()

	if annotations == nil {
		return obj, nil
	}
	if labels == nil {
		labels = map[string]string{}
	}

	if _, ok := obj.GetAnnotations()[extendedMonitoringAnnotationKey]; ok {
		delete(annotations, extendedMonitoringAnnotationKey)

		labels[extendedMonitoringAnnotationKey] = ""
	}

	obj.SetAnnotations(annotations)
	obj.SetLabels(labels)

	return obj, nil
}
