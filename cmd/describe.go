package cmd

import (
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return requireResource(cmd, args)
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
