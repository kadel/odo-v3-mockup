package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing components in the cluster.",
	Long: `List all component in the cluster in a given namespace.
By default it uses namespace as set by kubectl config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		namespace := "mynamespace"

		fmt.Printf("Components in the %q namespace:\n", namespace)

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New(" ", "NAME", "TYPE", "MANAGED BY ", "RUNNING IN").WithWriter(os.Stdout).WithPadding(2).WithWidthFunc(runewidth.StringWidth)
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		tbl.AddRow("*", "frontend", "nodejs", "odo", "Dev,Deploy")
		tbl.AddRow(" ", "backend", "springboot", "odo", "Dev")
		tbl.AddRow(" ", "created-by-odc", "python", "unknown", "unknown")
		tbl.AddRow(" ", "MariaDB", "unknown", "Helm", "unknown")
		tbl.Print()
		fmt.Println()

		fmt.Printf("ServiceBindings in the %q namespace:\n", namespace)
		tblSB := table.New("NAME", "Application", "SERVICES ").WithWriter(os.Stdout).WithPadding(2).WithWidthFunc(runewidth.StringWidth)
		tblSB.WithHeaderFormatter(headerFmt)
		tblSB.AddRow("backend-mongodb", "backend (Deploment)", "mymongodb (PrconaServerMongoDB.psmdb.percona.com)")
		tblSB.AddRow("frontend-redis", "frontend (Deployment)", "myredis (Redis.redis.redis.opstreelabs.in)")
		tblSB.AddRow("otherbinding", "application (Deployment)", "myredis (Redis.redis.redis.opstreelabs.in)")
		tblSB.Print()
		return nil
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
