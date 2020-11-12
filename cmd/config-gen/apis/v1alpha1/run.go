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
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

// NewCommand returns a new cobra command
func NewCommand() cobra.Command {
	fc := &APIConfiguration{}
	return framework.TemplateCommand{
		MergeResources: true,
		API:            fc,
		Template:       configTemplate,
		PreProcess: func(rl *framework.ResourceList) error {
			// Validate the input
			if fc.Name == "" {
				return errors.Errorf("must specify name")
			}
			if fc.Image == "" {
				return errors.Errorf("must specify image")
			}

			// Generate resources from the code
			var err error
			rl.Items, err = ControllerGenFilter{APIConfiguration: fc}.Filter(rl.Items)
			return err
		},
		PostProcess: func(rl *framework.ResourceList) error {
			// Sort the resources
			var err error
			rl.Items, err = SortFilter{APIConfiguration: fc}.Filter(rl.Items)
			return err
		},
	}.GetCommand()
}
