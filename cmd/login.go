package cmd

import (
	"github.com/kadel/odo-v3-prototype/pretty"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [URL]",
	Short: "Log in to your server and save login for subsequent use.",
	Long: `First-time users of the client should run this command to connect to a server, establish an authenticated session, and
save connection to the configuration file. The default configuration will be saved to your home directory under
".kube/config".

 The information required to login -- like username and password, a session token, or the server details -- can be
provided through flags. If not provided, the command will prompt for user input as needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		pretty.Printf(pretty.Success, "Login successful.")
		pretty.Printf(pretty.Info, "You are logged as %q on %q", "user", "https://crc.testing:6443")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("certificate-authority", "", "", "Path to a cert file for the certificate authority")
	loginCmd.Flags().StringP("password", "p", "", "Password for server")
	loginCmd.Flags().StringP("user", "u", "", "Username for server")
	loginCmd.Flags().StringP("token", "", "", "Bearer token for authentication to the API server")
	loginCmd.Flags().BoolP("insecure-skip-tls-verify", "", false, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
}
