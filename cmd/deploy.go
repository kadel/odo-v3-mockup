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
	"fmt"
	"os"
	"time"

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
		fmt.Println()
		color.New(color.Bold).Println("Building container image locally ...")
		fmt.Println("[1/2] STEP 1/3: FROM registry.access.redhat.com/ubi8/nodejs-14:latest")
		fmt.Println("[1/2] STEP 2/3: COPY package*.json ./")
		time.Sleep(time.Second * 1)
		fmt.Println("--> Using cache 1fe872faa61d81c61f01429364f035ecba3b553e96ce0d4582d1c0e61e1a171d")
		fmt.Println("--> 1fe872faa61")
		fmt.Println("[1/2] STEP 3/3: RUN npm install --production")
		time.Sleep(time.Second * 1)
		fmt.Println("--> Using cache bfa0e23103d476c30772f2b9f4fad07c4897c9dca46753fe85a5cb4d340e57d3")
		fmt.Println("--> bfa0e23103d")
		fmt.Println("[2/2] STEP 1/6: FROM registry.access.redhat.com/ubi8/nodejs-14-minimal:latest")
		fmt.Println("[2/2] STEP 2/6: COPY --from=0 /opt/app-root/src/node_modules /opt/app-root/src/node_modules")
		fmt.Println("--> Using cache 9bacf4fd410a999f064d5859a00d259f2a570c0313bced08984370058ee42b6e")
		fmt.Println("--> 9bacf4fd410")
		fmt.Println("[2/2] STEP 3/6: COPY . /opt/app-root/src")
		time.Sleep(time.Second * 1)
		fmt.Println("--> d5feee229cc")
		fmt.Println("[2/2] STEP 4/6: ENV NODE_ENV production")
		fmt.Println("--> e7bad84799f")
		fmt.Println("[2/2] STEP 5/6: ENV PORT 3000")
		fmt.Println("--> 1298e349f6d")
		fmt.Println("[2/2] STEP 6/6: CMD [\"npm\", \"start\"]")
		time.Sleep(time.Second * 1)
		fmt.Println("[2/2] COMMIT quay.io/tkral/devfile-nodejs-deploy:latest")
		fmt.Println("--> eba665280b6")
		fmt.Println("Successfully tagged quay.io/tkral/devfile-nodejs-deploy:latest")
		fmt.Println("eba665280b6b61327658620fbcdd57718cee1adea3615ad14e2ffd38fe5d33ca")
		fmt.Println()
		color.New(color.Bold).Println("Pushing image to container registry ...")
		fmt.Println("Getting image source signatures")
		time.Sleep(time.Second * 1)
		fmt.Println("Copying blob cd27cf2ab079 done  ")
		fmt.Println("Copying blob 54e42005468d skipped: already exists  ")
		fmt.Println("Copying blob ce97028d0f33 skipped: already exists  ")
		time.Sleep(time.Second * 1)
		fmt.Println("Copying blob 8353b13febf7 skipped: already exists  ")
		fmt.Println("Copying blob 0b911edbb97f skipped: already exists  ")
		time.Sleep(time.Second * 1)
		fmt.Println("Copying config eba665280b done  ")
		fmt.Println("Writing manifest to image destination")
		fmt.Println("Storing signatures")
		fmt.Println()
		Spinner("Waiting for Kubernetes resources ...", 3) // creates everything on the cluster and waits for contaiers to be ready
		color.Magenta("Your component is running on cluster. ")
		color.New(color.FgMagenta).Print("You can access it at ")
		color.New(color.FgMagenta).Add(color.Underline).Println("https://example.com")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
