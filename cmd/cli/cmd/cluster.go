package cmd

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/util/completion"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	clusterCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default",
		"Decide how you want data to be displayed (json,yaml,wide,name,clusterid,environment,datacenter,lastObserved,firstObserved,created)")
	clusterCmd.PersistentFlags().BoolVarP(&envFlag, "env", "e", false, "Show the environment of the current cluster\nIts a lot faster than:\n $ kubectl -n argocd get applications nhn-tooling -o json | jq .spec.source.helm.values -r | \\ \n   awk -F'\\n' '{ print $1 \"\\n\" $2 }' | yq -o json | jq .nhn.environment | sed 's/\"//g'")
	rootCmd.AddCommand(clusterCmd)
}

var envFlag bool
var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"cl"},
	Short:   "Get details about a specific cluster",
	Long: `
	Display a cluster you have access too based on its name

	Example:

	ror cluster tooling-staging.trd1-nhn-tooling

	`,
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completion.ClusterAndWorkspaceIdCompletionFunc,
	PreRun:            cmdrorclient.SetupRorClient,
	Run:               cmdCluster,
}

func cmdCluster(cmd *cobra.Command, args []string) {
	var id string

	if len(args) > 0 {
		id = args[0]
	} else {
		id = viper.GetString(config.LastSessionCluster)
	}

	if id == "" {
		rlog.Fatal("error executing command", fmt.Errorf("missing a clusterId, could not find it in arguments or lastsession"))
	}

	if envFlag {
		cluster, err := config.RorClient.Clusters().GetById(id)
		cobra.CheckErr(err)
		fmt.Println(cluster.Environment)
		return
	}

	cluster, err := config.RorClient.Clusters().GetById(id)
	cobra.CheckErr(err)

	var clusters [1]apicontracts.Cluster
	clusters[0] = *cluster
	if !cmd.Flags().Changed("output") {
		outputformatting.ClusterTabPrinter(clusters[:])
		return
	}

	outputformatting.HandleClusterOutputFormatting(outputFormat, clusters[:])
}
