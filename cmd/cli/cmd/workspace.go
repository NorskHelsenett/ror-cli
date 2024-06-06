package cmd

import (
	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/util/completion"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	workspaceCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Decide how you want data to be displayed (json/default)")
	rootCmd.AddCommand(workspaceCmd)
}

var workspaceCmd = &cobra.Command{
	Use:     `workspace <workspace-name>`,
	Aliases: []string{"wo"},
	Short:   "Get details about a specific workspace",
	Long: `
	Display information about a workspace you have access to based on its name.
	The name is either infered by the last cluster you logged into or by a argument

	Example:

	ror workspace trd1-nhn-tooling

	`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completion.WorkspaceIdCompletionFunc(cmd, args, toComplete)
	},
	PreRun: cmdrorclient.SetupRorClient,
	Run:    cmdWorkspace,
}

func cmdWorkspace(cmd *cobra.Command, args []string) {
	var name string

	if len(args) > 0 {
		name = args[0]
	} else {
		name = viper.GetString(config.LastSessionWorkspace)
	}

	if name == "" {
		_ = cmd.Help()
	}

	workspace, err := config.RorClient.Workspaces().GetByName(name)
	cobra.CheckErr(err)

	var workspaces [1]apicontracts.Workspace
	workspaces[0] = *workspace
	if !cmd.Flags().Changed("output") {
		outputformatting.WorkspaceTabPrinter(workspaces[:])
		return
	}

	outputformatting.HandleWorkspaceOutputFormatting(outputFormat, workspaces[:])
}
