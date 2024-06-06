package cmd

import (
	"fmt"
	"os"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/cobra"
)

func init() {
	workspacesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Decide how you want data to be displayed (json/default)")
	rootCmd.AddCommand(workspacesCmd)
}

var workspacesCmd = &cobra.Command{
	Use:     "workspaces",
	Aliases: []string{"w"},
	Short:   "Get details about all workspaces",
	Long: `
	Display information about all the workspaces you have access to.

	Example:

	ror workspaces

	`,
	PreRun: cmdrorclient.SetupRorClient,
	Run:    cmdWorkspaces,
}

func cmdWorkspaces(cmd *cobra.Command, args []string) {
	workspaces, err := config.RorClient.Workspaces().GetAll()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "could not get workspaces from ror api")
		rlog.Fatal("could not get workspaces from ror api: ", err)
	}

	if !cmd.Flags().Changed("output") {
		outputformatting.WorkspaceTabPrinter(*workspaces)
		return
	}

	outputformatting.HandleWorkspaceOutputFormatting(outputFormat, *workspaces)
}
