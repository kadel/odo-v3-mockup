package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var createService = &cobra.Command{
	Use:   "service",
	Short: "Create new service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {

		if !HasFlagsSet(cmd) {
			serviceQuesion := &survey.Select{
				Message: "Select Service:",
				Options: services,
			}
			var serviceAnswer string
			survey.AskOne(serviceQuesion, &serviceAnswer)

			serviceNameQuestion := &survey.Input{Message: "What will be new service name?"}
			var serviceNameAnswer string
			survey.AskOne(serviceNameQuestion, &serviceNameAnswer)

		}
	},
}

var deleteService = &cobra.Command{
	Use:   "service",
	Short: "Delete existing service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")

	},
}

var listService = &cobra.Command{
	Use:   "service",
	Short: "List existing services",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")

	},
}

func init() {
	createCmd.AddCommand(createService)
	deleteCmd.AddCommand(deleteService)
	listCmd.AddCommand(listService)
}
