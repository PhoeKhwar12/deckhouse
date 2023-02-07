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
	"github.com/deckhouse/deckhouse/go_lib/hooks/tls_certificate"
)

const (
	cn = "extended-monitoring-webhook"
)

var _ = tls_certificate.RegisterInternalTLSHook(tls_certificate.GenSelfSignedTLSHookConf{
	SANs: tls_certificate.DefaultSANs([]string{
		"extended-monitoring-webhook.d8-monitoring.svc",
		"annotation-converter-webhook.d8-monitoring.svc",
		"extended-monitoring-webhook.d8-monitoring",
		"annotation-converter-webhook.d8-monitoring",
		"extended-monitoring-webhook",
		"annotation-converter-webhook",
	}),

	CN:                   cn,
	Namespace:            "d8-monitoring",
	TLSSecretName:        "extended-monitoring-webhook-tls",
	FullValuesPathPrefix: "extendedMonitoring.internal.webhookCert",
})
