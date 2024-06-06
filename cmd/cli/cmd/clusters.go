package cmd

import (
	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"

	"github.com/spf13/cobra"
)

func init() {
	clustersCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default",
		"Decide how you want data to be displayed (json,yaml,wide,name,clusterid,environment,datacenter,lastObserved,firstObserved,created)")
	clustersCmd.Flags().IntVar(&paginationLimit, "limit", 50, "Paginated amount of clusters.")
	clustersCmd.Flags().IntVar(&paginationOffset, "offset", 0, "Offset for pagination.")

	rootCmd.AddCommand(clustersCmd)
}

var clustersCmd = &cobra.Command{
	Use:     "clusters",
	Aliases: []string{"c"},
	Short:   "Get details about all clusters",
	Long: `
	Display all the clusters you have access to.

	Example:

	ror clusters
	`,
	PreRun: cmdrorclient.SetupRorClient,
	Run:    cmdClusters,
}

func cmdClusters(cmd *cobra.Command, args []string) {
	clusters, err := config.RorClient.Clusters().GetAll()
	cobra.CheckErr(err)

	if !cmd.Flags().Changed("output") {
		outputformatting.ClusterTabPrinter(*clusters)
		return
	}
	outputformatting.HandleClusterOutputFormatting(outputFormat, *clusters)
}
