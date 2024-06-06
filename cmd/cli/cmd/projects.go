package cmd

import (
	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"

	"github.com/spf13/cobra"
)

func init() {
	projectsCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default",
		"Decide how you want data to be displayed (json,yaml,wide,name,clusterid,environment,datacenter,lastObserved,firstObserved,created)")

	rootCmd.AddCommand(projectsCmd)
}

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"p"},
	Short:   "Get details about all projects",
	Long: `
	Displays all projects you have access to

	Example:

	ror projects

	`,

	Args:   cobra.MaximumNArgs(1),
	PreRun: cmdrorclient.SetupRorClient,
	Run:    cmdProjects,
}

func cmdProjects(cmd *cobra.Command, args []string) {
	projects, err := config.RorClient.Projects().GetAll()
	cobra.CheckErr(err)
	if !cmd.Flags().Changed("output") {
		outputformatting.ProjectTabPrinter(*projects)
		return
	}

	outputformatting.HandleProjectOutputFormatting(outputFormat, *projects)
}
