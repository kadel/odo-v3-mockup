package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var createEndpoint = &cobra.Command{
	Use:   "endpoint",
	Short: "Create new endpoint",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")
	},
}

var deleteEndpoint = &cobra.Command{
	Use:   "endpoint",
	Short: "Delete existing endpoint",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")
	},
}

var listEndpoint = &cobra.Command{
	Use:   "endpoint",
	Short: "List existing endpoints",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")
	},
}

func init() {
	createCmd.AddCommand(createEndpoint)
	deleteCmd.AddCommand(deleteEndpoint)
	listCmd.AddCommand(listEndpoint)
}
