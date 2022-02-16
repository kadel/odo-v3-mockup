package cmd

import (
	"github.com/kadel/odo-v3-prototype/pretty"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of the active session out by clearing saved tokens.",
	Long: `An authentication token is stored in the config file after login - this command will delete that token on the server,
and then remove the token from the configuration file.

 If you are using an alternative authentication method like Kerberos or client certificates, your ticket or client
certificate will not be removed from the current system since these are typically managed by other programs. Instead,
you can delete your config file to remove the local copy of that certificate or the record of your server login.

 After logging out, if you want to log back into the server use 'oc login'.`,
	Run: func(cmd *cobra.Command, args []string) {
		pretty.Printf(pretty.Success, "Logout successfully from %q", "https://crc.testing:6443")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

}
