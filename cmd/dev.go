package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run application in a developer mode.",
	Long: `Run application in a developer mode.
Application will be started in a development mode on the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("devfile.yaml"); os.IsNotExist(err) {
			color.Green("There is no devfile.yaml in the current directory.")
			devfileName, devfileRegistry, projectName := SelectDevfileAlizer(cmd)
			DownloadDevfile(devfileName, devfileRegistry, projectName)
		} else {
			color.Green("Using devfile.yaml from the current directory.")
		}
		color.Green("Starting your application on cluster in developer mode ...")
		Spinner("Waiting for Kubernetes resources ...", 2) // creates everything on the cluster and waits for contaiers to be ready
		Spinner("Syncing files into the container ...", 2) // push the files into the container
		Spinner("Building your application in container on cluster ...", 2)         // start the application
		Spinner("Execting the application ...", 2)         // start the application
		color.Magenta("Your application is running on cluster. ")
		color.New(color.FgMagenta).Print("You can access it at ")
		color.New(color.FgMagenta).Add(color.Underline).Println("https://example.com")
	},
}

func init() {
	addCommonDevfileFlags(devCmd)
	rootCmd.AddCommand(devCmd)
}
