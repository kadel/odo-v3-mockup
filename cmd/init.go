package cmd

import (
	"fmt"

	"github.com/devfile/registry-support/index/generator/schema"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap new application based on Devfile",
	Long:  `Bootstrap new application based on Devfile`,
	Run: func(cmd *cobra.Command, args []string) {

		var devfile schema.Schema
		var devfileRegistry string
		var componentName string
		var starterName string

		if !IsCurrentDirEmpty() {
			devfile, devfileRegistry, componentName = SelectDevfileAlizer(cmd)
		} else {
			devfile, devfileRegistry, componentName, starterName = SelectDevfile(cmd, true)
		}

		DownloadDevfile(devfile, devfileRegistry, componentName, starterName)

		color.Green("Your new component %q is ready in the current directory.\n", componentName)
		fmt.Println("To start editing your project, use “odo dev” and open this folder in your favorite IDE.")
		fmt.Println("Changes will be directly reflected on the cluster.")
		fmt.Println("To deploy your application to your cluster use “odo deploy”.")
	},
}

func init() {
	initCmd.Flags().String("devfile-path", "", "Path to a devfile. It can be local filesystem path or http(s). This is alternative to using devfile from Devfile registry.")
	initCmd.Flags().String("devfile", "", "Name of the Devfile from the Devfile registry (required if --devfile-path is not defined)")
	initCmd.Flags().String("registry", "", "name of the devfile registry (as configured in odo preference registry). It can be used in combination with --devfile, but not with --devfile-paths")
	initCmd.Flags().String("name", "", "Name of the new component. If not specified, the name will be the name of the current directory.")
	initCmd.Flags().String("starter", "", "Name of the devfile starter project to populate the current directory with.")

	rootCmd.AddCommand(initCmd)
}
