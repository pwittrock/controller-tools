/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// APIConfiguration implements the API for generating configuration
type APIConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Name is the name of the config
	Name string `json:"name" yaml:"name"`

	// Image is the container image to use in the controller-manager
	Image string `json:"image" yaml:"image"`

	// DisableCreateNamespace prevents the namespace from being generated
	DisableCreateNamespace bool `json:"disableCreateNamespace,omitempty" yaml:"disableCreateNamespace,omitempty"`

	// DisableCreateManager if set to true will disable generating the controller-manager
	DisableCreateManager bool `json:"disableCreateManager,omitempty" yaml:"disableCreateManager,omitempty"`

	// DisableCreateRBAC if set to true will disable generating the rbac
	DisableCreateRBAC bool `json:"disableCreateRBAC,omitempty" yaml:"disableCreateRBAC,omitempty"`

	// DisableAuthProxy if set to true will disable the auth proxy
	DisableAuthProxy bool `json:"disableAuthProxy,omitempty" yaml:"disableAuthProxy,omitempty"`

	// EnableWebhooks configures webhooks for the controller-manager
	EnableWebhooks bool `json:"enableWebhooks,omitempty" yaml:"enableWebhooks,omitempty"`

	// EnableCertManager uses the jetstack certmanager to inject certificates
	// for webhooks.
	EnableCertManager bool `json:"enableCertManager,omitempty" yaml:"enableCertManager,omitempty"`

	// EnablePrometheus creates a service monitor
	EnablePrometheus bool `json:"enablePrometheus,omitempty"`

	// ConversionWebhooks is a map of kinds to enable conversion webhooks for
	ConversionWebhooks map[string]bool `json:"conversionWebhooks,omitempty"`
}
