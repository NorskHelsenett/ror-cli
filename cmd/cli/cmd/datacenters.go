package cmd

import (
	"fmt"
	"os"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"

	"github.com/spf13/cobra"
)

func init() {
	datacentersCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Decide how you want data to be displayed (json/default)")
	rootCmd.AddCommand(datacentersCmd)
}

var datacentersCmd = &cobra.Command{
	Use:     "datacenters",
	Aliases: []string{"d"},
	Short:   "Get details about all datacenters",
	Long:    `Display all the datacenters you have access to.`,
	PreRun:  cmdrorclient.SetupRorClient,
	Run:     cmdDatacenters,
}

func cmdDatacenters(cmd *cobra.Command, args []string) {
	datacenters, err := config.RorClient.Datacenters().Get()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "could not get datacenter from ror api")
	}

	if !cmd.Flags().Changed("output") {
		outputformatting.DatacenterTabPrinter(*datacenters)
		return
	}

	outputformatting.HandleDatacenterOutputFormatting(outputFormat, datacenters)
}
