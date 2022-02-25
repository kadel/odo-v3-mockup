package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources from the cluster",
	//Long:  ``,
	RunE: requireResource,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(deleteCmd)

}

var deleteComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Delete component from the cluster.",
	Long: `If the command is exected in a directory with 'devfile.yaml'
it will try to delete all cluster resources associated with the component.	
Optionaly, you can use '--name' and '--namespace' flags to delete a specific component.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		var namespace string

		//TODO: this should ve read from local devfile
		name = "component-from-devfile"
		//TODO: this should be read form .env
		namespace = "mynamespace"

		if cmd.Flag("name").Changed {
			name = cmd.Flag("name").Value.String()
		}

		if cmd.Flag("namespace").Changed {
			name = cmd.Flag("namespace").Value.String()
		}

		var confirmAnswer bool
		confirm := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to delete %q component from %q namespace?", name, namespace),
			Default: true,
		}
		survey.AskOne(confirm, &confirmAnswer)

		if confirmAnswer {
			Spinner(fmt.Sprintf("Deleting component %q from %q namespace", name, namespace), 3)
		}
		return nil
	},

	Args: cobra.NoArgs,
}

func init() {
	deleteComponentCmd.Flags().String("name", "", "Name of the new component that will be deleted.")
	deleteComponentCmd.Flags().String("namespace", "", "Namespace from which the component will be deleted.")
	deleteComponentCmd.Flags().Bool("force", false, "Don't ask for confirmation.")
	deleteCmd.AddCommand(deleteComponentCmd)

}
