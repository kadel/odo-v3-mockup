package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func requireResource(cmd *cobra.Command, args []string) error {
	var subCommands string
	for _, cmd := range cmd.Commands() {
		subCommands += cmd.Name() + " "
	}
	return fmt.Errorf("you need to specify resource to create. Available resources are: %s", subCommands)
}

// HasFlagsSet returns true if any of the flags in the given command is set
func HasFlagsSet(cmd *cobra.Command) bool {
	hasFlagsSet := false
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			hasFlagsSet = true
		}
	})

	return hasFlagsSet
}

// addCommonFlags adds flags that are used to download devfile from registry
func addCommonDevfileFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Your project's name")
	cmd.Flags().String("devfile", "", "Name of the devfile in Devfile registry")
	cmd.Flags().String("registry", "defaultRegistry", "Name of the devfile registry to download devfile from")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
