// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/scaffold"
	"sigs.k8s.io/controller-tools/pkg/scaffold/manager"
	"sigs.k8s.io/controller-tools/pkg/scaffold/project"
)

var prj *project.Project
var bp *project.Boilerplate
var gopkg *project.GopkgToml
var mrg *manager.Cmd
var dkr *manager.Dockerfile

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Scaffold a basic project.",
	Long: `Scaffold a basic project including:

- Writing or updating a Gopkg.toml with project dependencies
- Writing a boilerplate license file
- Writing a PROJECT file with the domain
- Writing a cmd/manager/main.go to run Controllers
`,
	Example: `# Scaffold a project
controller-tools scaffold project --domain k8s.io --license apache2 --owner "The Kubernetes authors"

# Fetch the dependencies
dep ensure
`,
	Run: func(cmd *cobra.Command, args []string) {
		// project and boilerplate must come before main so the boilerplate exists
		s := &scaffold.Scaffold{
			BoilerplateOptional: true,
			ProjectOptional:     true,
		}
		err := s.Execute(scaffold.Options{
			ProjectPath:     prj.Path(),
			BoilerplatePath: bp.Path(),
		}, prj, bp)
		if err != nil {
			log.Fatal(err)
		}

		s = &scaffold.Scaffold{}
		err = s.Execute(scaffold.Options{
			ProjectPath:     prj.Path(),
			BoilerplatePath: bp.Path(),
		}, gopkg, mrg, dkr, &manager.APIs{}, &manager.Controller{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Run `dep ensure` to fetch dependencies [y/n]?")
		if yesno() {
			c := exec.Command("dep", "ensure")
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			fmt.Println(strings.Join(c.Args, " "))
			if err := c.Run(); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Skipping `dep ensure`.  Dependencies will not be fetched.")
		}

		fmt.Println("Next: Scaffold a new API with `controller-tools scaffold resource`.")
	},
}

func init() {
	scaffoldCmd.AddCommand(projectCmd)

	prj = project.ForFlags(projectCmd.Flags())
	bp = project.BoilerplateForFlags(projectCmd.Flags())
	gopkg = project.GopkgTomlForFlags(projectCmd.Flags())
	mrg = manager.ForFlags(projectCmd.Flags())
	dkr = manager.DockerfileForFlags(projectCmd.Flags())
}
