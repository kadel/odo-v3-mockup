package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing resources",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	RunE: requireResource,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
