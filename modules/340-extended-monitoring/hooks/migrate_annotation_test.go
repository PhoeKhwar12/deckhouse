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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const annotatedObjects = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
  namespace: default
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
  namespace: default
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
  namespace: default
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
  namespace: default
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
  namespace: default
---
apiVersion: v1
kind: Node
metadata:
  annotations:
    extended-monitoring.flant.com/enabled: ""
  name: test
`

var _ = Describe("Modules :: extended-monitoring :: hooks :: migrate_annotations ::", func() {

	f := HookExecutionConfigInit(`{}`, `{}`)

	Context("With annotated objects", func() {
		BeforeEach(func() {
			f.KubeStateSet(annotatedObjects)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		It("All extended-monitoring annotations should be replaced with labels", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.PatchCollector.Operations()).To(HaveLen(6))

			deployment := f.KubernetesResource("Deployment", "default", "test")
			Expect(deployment.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(deployment.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())

			statefulSet := f.KubernetesResource("StatefulSet", "default", "test")
			Expect(statefulSet.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(statefulSet.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())

			daemonSet := f.KubernetesResource("DaemonSet", "default", "test")
			Expect(daemonSet.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(daemonSet.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())

			cronJob := f.KubernetesResource("CronJob", "default", "test")
			Expect(cronJob.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(cronJob.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())

			ingress := f.KubernetesResource("Ingress", "default", "test")
			Expect(ingress.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(ingress.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())

			node := f.KubernetesResource("Node", "", "test")
			Expect(node.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeFalse())
			Expect(node.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeTrue())
		})
	})

	Context("Objects with legacy annotation already migrated", func() {
		BeforeEach(func() {
			f.KubeStateSet(`---
apiVersion: v1
kind: ConfigMap
metadata:
  name: d8-extended-monitoring-migrated
  namespace: d8-monitoring
` + annotatedObjects)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		It("Object annotations should not change", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.PatchCollector.Operations()).To(HaveLen(0))

			deployment := f.KubernetesResource("Deployment", "default", "test")
			Expect(deployment.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(deployment.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())

			statefulSet := f.KubernetesResource("StatefulSet", "default", "test")
			Expect(statefulSet.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(statefulSet.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())

			daemonSet := f.KubernetesResource("DaemonSet", "default", "test")
			Expect(daemonSet.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(daemonSet.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())

			cronJob := f.KubernetesResource("CronJob", "default", "test")
			Expect(cronJob.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(cronJob.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())

			ingress := f.KubernetesResource("Ingress", "default", "test")
			Expect(ingress.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(ingress.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())

			node := f.KubernetesResource("Node", "", "test")
			Expect(node.Field(`metadata.annotations.extended-monitoring\.flant\.com/enabled`).Exists()).To(BeTrue())
			Expect(node.Field("metadata.labels.extended-monitoring\\.flant\\.com/enabled").Exists()).To(BeFalse())
		})
	})
})
