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

package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	whhttp "github.com/slok/kubewebhook/v2/pkg/http"
	"github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
)

const extendedMonitoringAnnotationKey = "extended-monitoring.flant.com/enabled"

func mutateObject(_ context.Context, _ *kwhmodel.AdmissionReview, obj metav1.Object) (mr *mutating.MutatorResult, err error) {
	annotations := obj.GetAnnotations()
	labels := obj.GetLabels()

	if annotations == nil {
		return &mutating.MutatorResult{MutatedObject: obj}, nil
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

	return &mutating.MutatorResult{MutatedObject: obj}, nil
}

type config struct {
	certFile   string
	keyFile    string
	listenAddr string
	debug      bool
}

func initFlags() *config {
	cfg := &config{}

	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fl.StringVar(&cfg.certFile, "tls-cert-file", "", "TLS certificate file")
	fl.StringVar(&cfg.keyFile, "tls-key-file", "", "TLS key file")
	fl.StringVar(&cfg.listenAddr, "listen-address", ":8080", "listen address")
	fl.BoolVar(&cfg.debug, "debug", false, "debug logging")

	klog.InitFlags(fl)

	_ = fl.Parse(os.Args[1:])
	if cfg.certFile == "" && cfg.keyFile == "" {
		klog.Fatal(`"tls-cert-file" and/or "tls-key-file" args not provided`)
	}
	return cfg
}

func main() {
	cfg := initFlags()
	logrusLogEntry := logrus.NewEntry(logrus.New())
	logLevel := logrus.WarnLevel
	if cfg.debug {
		logLevel = logrus.DebugLevel
	}
	logrusLogEntry.Logger.SetLevel(logLevel)
	kl := kwhlogrus.NewLogrus(logrusLogEntry)

	wh, err := mutating.NewWebhook(
		mutating.WebhookConfig{
			ID:      "convertAnnotation",
			Obj:     &unstructured.Unstructured{},
			Mutator: mutating.MutatorFunc(mutateObject),
			Logger:  kl,
		})
	if err != nil {
		klog.Fatalf("error creating webhook: %s", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/mutate", whhttp.MustHandlerFor(whhttp.HandlerConfig{Webhook: wh, Logger: kl}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })

	klog.Info("Listening on :8080")

	err = http.ListenAndServeTLS(cfg.listenAddr, cfg.certFile, cfg.keyFile, mux)
	if err != nil {
		klog.Fatalf("error serving webhook: %s", err)
	}
}
