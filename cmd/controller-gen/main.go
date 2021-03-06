/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"
	"os"
	"strings"

	"sigs.k8s.io/controller-tools/cmd/controller-gen/cmd"
)

//go:generate go run ../helpgen/main.go paths=../../pkg/... generate:headerFile=../../boilerplate.go.txt,year=2019

func main() {
	c := cmd.New()

	if err := c.Execute(); err != nil {
		if _, noUsage := err.(cmd.NoUsageError); !noUsage {
			// print the usage unless we suppressed it
			if err := c.Usage(); err != nil {
				panic(err)
			}
		}
		fmt.Fprintf(c.OutOrStderr(), "run `%[1]s %[2]s -w` to see all available markers, or `%[1]s %[2]s -h` for usage\n", c.CalledAs(), strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}
}
