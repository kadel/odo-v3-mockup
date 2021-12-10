package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap new application based on Devfile",
	Long:  `Bootstrap new application based on Devfile`,
	Run: func(cmd *cobra.Command, args []string) {

		if !IsCurrentDirEmpty() {
			color.Red("Current directory is not empty. You can bootstrap new application only in empty directory.")
			color.Green("If you have existing code that you want to deploy use `%s deploy` or use `%s dev` command to quickly iterate on your application.", os.Args[0], os.Args[0])
			os.Exit(1)
		}

		//devfile, devfileRegistry, projectName := SelectDevfileFromRegistry(cmd)
		devfile, devfileRegistry, componentName, starterName := SelectDevfile(cmd, true)

		DownloadDevfile(devfile, devfileRegistry, componentName, starterName)

		color.Green("Your new component %q is ready in the current directory.\n", componentName)
		fmt.Println("To start editing your project, use “odo dev” and open this folder in your favorite IDE.")
		fmt.Println("Changes will be directly reflected on the cluster.")
		fmt.Println("To deploy your application to your cluster use “odo deploy”.")
	},
}

func init() {
	addCommonDevfileFlags(initCmd)
	initCmd.Flags().String("starter", "", "Name of the devfile starter project to populate the current directory with")

	rootCmd.AddCommand(initCmd)
}
