package cmd

import (
	"github.com/NorskHelsenett/ror-cli/cmd/cli/cmd/login"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/util/completion"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	loginCmd.Flags().BoolP("workspace", "w", false, "Tell the command to use workspaces instead of clusters when presenting the list or logging directly into a cluster with \"ror login -w trd1-nhn-tooling\"")

	logoutCmd.Flags().BoolP("ror", "r", false, "Reset ror user data")
	logoutCmd.Flags().BoolP("privileged", "p", false, "Reset privileged user data")
}

var loginCmd = &cobra.Command{
	Use:     "login {ClusterId/WorkspaceName}",
	Aliases: []string{"l"},
	Short:   `Multi-use command that authenicates you and logs you into a cluster or worksapce`,
	Long: `
	A command that can be used to log into a kubernetes cluster or a workspace indexed by ROR, 
	the command will try to sign you in to a cluster based on its provider.

	When first using the command you will be prompted for authentication, this is done by
	opening a browser window and asking you to log in to your identity provider (Single Sign On).

	Then the program will try to connect to PAM (Cyberark) to get your privileged credentials.
	Tou will be prompted for your PAM password, this is the password you use to log into the pam portal.
	This will result in a 2FA flow, where you will be prompted for  verification in your 2FA app.

	If ror-cli is unable to connect to PAM, you will be prompted for your privileged user and password.

	ror-cli will then request the kube-context and merge it with your current kube-config.
	
	If you alreday know the ClusterId if the cluster you want to log into you can pass it as a argument
	too the command. "ror login trd1-nhn-tooling-xxxx" for example, here you can also be prompted for 
	authentication. 

	Specifing that we want to log into workspaces or a workspace is as easy as using the "-w" flag.	
	Here you would also be prompted with a list of workspaces instead of clusters, which you can log into.

	It's also possible to tell "ror" to log into a workspace we know the WorkspaceId of.
	"ror login -w trd1-nhn-tooling"

	To logout of the current session use "ror logout"
	  ror logout --help for more information

	`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completion.ClusterAndWorkspaceIdCompletionFunc,
	PreRun:            cmdrorclient.SetupRorClient,
	Run:               login.RorLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: `Log out the currently logged in user`,
	Long: `
	Logs out the logged in user, if no flags specified both privileged and ror userdata vil be removed.
	`,
	Args: cobra.MaximumNArgs(0),
	Run:  login.RorLogout,
}
