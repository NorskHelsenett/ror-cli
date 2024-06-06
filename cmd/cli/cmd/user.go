package cmd

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u"},
	Short:   "Show information about the current authenticated user including username, email, and groups",
	Long:    `List user information of the authenticated user,  including username, email, and groups`,
	PreRun:  cmdrorclient.SetupRorClient,
	Run:     printUserInfo,
}

func printUserInfo(cmd *cobra.Command, args []string) {
	fmt.Println("UserInfo")
	fmt.Println("--------")
	fmt.Println("User:", config.Authinfo.User.Email)
	fmt.Println("Name:", config.Authinfo.User.Name)
	fmt.Println("Groups:")
	for _, group := range config.Authinfo.User.Groups {
		fmt.Println("  -", group)
	}
	fmt.Println("")
	fmt.Println("AuthInfo")
	fmt.Println("--------")
	fmt.Println("Authenticated using:", config.Authinfo.Auth.AuthProvider)
	fmt.Println("Expires:", config.Authinfo.Auth.ExpirationTime.Local().Format("15:04 02.01.2006 MST"))

}
