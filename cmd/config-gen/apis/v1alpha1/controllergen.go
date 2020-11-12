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
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var _ kio.Filter = &ControllerGenFilter{}

// ControllerGenFilter generates resources using controller-gen
type ControllerGenFilter struct {
	*APIConfiguration
}

// Filter implements kio.Filter
func (cgr ControllerGenFilter) Filter(input []*yaml.RNode) ([]*yaml.RNode, error) {
	c := exec.Command("controller-gen", "paths=./...", "crd", "output:crd:stdout")
	if !cgr.DisableCreateRBAC {
		c.Args = append(c.Args, "output:rbac:stdout", fmt.Sprintf("rbac:roleName=%s-manager-role", cgr.Name))
	}
	if cgr.EnableWebhooks {
		c.Args = append(c.Args, "webhook", "output:webhook:stdout")
	}

	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = os.Stderr

	// Generate resources from the code and use them as input
	if err := c.Run(); err != nil {
		return nil, errors.WrapPrefixf(err, "failed to run controller-gen")
	}

	// Parse the output
	n, err := (&kio.ByteReader{Reader: &out}).Read()
	if err != nil {
		return nil, errors.WrapPrefixf(err, "failed to parse controller-gen output")
	}
	return append(n, input...), nil
}
