package completion

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"

	"github.com/spf13/cobra"
)

func ClusterAndWorkspaceIdCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cmdrorclient.SetupRorClient(cmd, args)
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var comp []string
	directive := cobra.ShellCompDirectiveNoFileComp

	isWorkspace, _ := cmd.Flags().GetBool("workspace")

	if !isWorkspace {
		comp = clusterIdComp(toComplete)
	}
	if isWorkspace {
		comp = workspaceIdComp(toComplete)
	}

	return comp, directive
}

func WorkspaceIdCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cmdrorclient.SetupRorClient(cmd, args)
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	directive := cobra.ShellCompDirectiveNoFileComp

	comp := workspaceIdComp(toComplete)

	return comp, directive
}

func clusterIdComp(toComplete string) []string {
	var clusterIdComp []string
	field := "clusterid"

	config.Load("")

	sortMetadata := apicontracts.SortMetadata{
		SortField: field,
		SortOrder: 1,
	}

	filterMetadata := apicontracts.FilterMetadata{
		Field:     field,
		Value:     toComplete,
		MatchMode: apicontracts.MatchModeContains,
	}

	filter := apicontracts.Filter{
		Skip:    0,
		Limit:   50,
		Sort:    []apicontracts.SortMetadata{sortMetadata},
		Filters: []apicontracts.FilterMetadata{filterMetadata},
	}
	clusters, err := config.RorClient.Clusters().GetByFilter(filter)
	if err != nil {
		msg := fmt.Sprintf("error from ror: %v\n", err)
		cobra.CompDebug(msg, true)
		return nil
	}

	for _, cluster := range *clusters {
		clusterIdComp = append(clusterIdComp, cluster.ClusterId)
	}

	return clusterIdComp
}

func workspaceIdComp(toComplete string) []string {
	var clusterIdComp []string

	config.Load("")

	workspaces, err := config.RorClient.Workspaces().GetAll()
	if err != nil {
		msg := fmt.Sprintf("error from ror: %v\n", err)
		cobra.CompDebug(msg, true)
		return nil
	}

	for _, workspace := range *workspaces {
		clusterIdComp = append(clusterIdComp, workspace.Name)
	}

	return clusterIdComp
}
