package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete  existing resource",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: requireResource,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}
