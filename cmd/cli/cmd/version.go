package cmd

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Shows you what version of ror you are running",
	Long:    `All software has versions. This is NHN-ROR CLI's`,
	PreRun:  cmdrorclient.SetupRorNonAuthClient,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {

	fmt.Println("Client: ", config.Version)
	serverversion := viper.GetString(config.ServerVersion)
	if serverversion == "" {
		serverversion = "Unknown/Unreachable"
	}
	fmt.Println("Server: ", serverversion)
}
