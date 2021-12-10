/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy application to Kubernetes cluster.",
	Long: `Deploy application to Kubernetes cluster.
It works in 3 steps:
- build the image
- push image to the container registry
- deploy it to Kubernetes cluster
`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("devfile.yaml"); os.IsNotExist(err) {
			color.Green("There is no devfile.yaml in the current directory.")
			devfile, devfileRegistry, projectName := SelectDevfileAlizer(cmd)
			DownloadDevfile(devfile, devfileRegistry, projectName, "")
		} else {
			color.Green("Using devfile.yaml from the current directory.")
		}
		color.Green("Deploying your component to cluster ...")
		Spinner("Building container image locally ...", 2)
		Spinner("Pushing image to container registry ...", 2)
		Spinner("Waiting for Kubernetes resources ...", 2) // creates everything on the cluster and waits for contaiers to be ready
		color.Magenta("Your component is running on cluster. ")
		color.New(color.FgMagenta).Print("You can access it at ")
		color.New(color.FgMagenta).Add(color.Underline).Println("https://example.com")
	},
}

func init() {
	addCommonDevfileFlags(deployCmd)
	rootCmd.AddCommand(deployCmd)
}
