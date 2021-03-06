/*
Copyright 2019 The Knative Authors
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
	"os"
	"strings"

	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"

	channelwebhook "knative.dev/eventing-kafka/pkg/channel/webhook"
)

const (
	Component = "kafkachannel-webhook"
	Secret    = "messaging-webhook-certs"
)

func main() {

	// Optionally Enable Support For ResetOffset
	if strings.ToLower(os.Getenv("RESETOFFSET_SUPPORT")) == "true" {
		channelwebhook.IncludeResetOffset()
	}

	// Define Webhook Options
	options := webhook.Options{
		ServiceName: webhook.NameFromEnv(),
		Port:        webhook.PortFromEnv(8443),
		SecretName:  Secret,
	}

	// Create A Signal Context With Webhook Options
	ctx := webhook.WithOptions(signals.NewContext(), options)

	// Create The Webhook With Desired Controllers
	sharedmain.MainWithContext(ctx, Component,
		certificates.NewController,
		channelwebhook.NewDefaultingAdmissionController,
		channelwebhook.NewValidationAdmissionController,
		channelwebhook.NewConversionController,
		// TODO(mattmoor): Support config validation in eventing-kafka.
	)
}
