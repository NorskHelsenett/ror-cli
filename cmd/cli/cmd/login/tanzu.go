package login

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"

	//"github.com/NorskHelsenett/ror-cli/cmd/cli/services/kubeconfig"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"
	"github.com/NorskHelsenett/ror/pkg/helpers/kubeconfig"
	"github.com/NorskHelsenett/ror/pkg/rlog"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func loginToTanzuCluster(cluster apicontracts.Cluster) error {
	var privilegedusername, privilegedpassword string
	// we dont have to login if the cluster token is still valid
	// Todo: use clusterid instead of clustername

	if cluster.Environment == "prod" {
		_, _ = fmt.Printf("\n   %s\n\n", color.HiRedString("Warning: You are logging into a production cluster"))
	}

	if !checkContextAndSwitch(cluster.ClusterName) {
		var err error
		privilegedusername, privilegedpassword, err = getPrivilegedCreds()
		cobra.CheckErr(err)

		fmt.Println("Logging into cluster " + colorizeByCluster(&cluster) + " with user " + color.YellowString(privilegedusername) + "")
		kubereturn, err := config.RorClient.Clusters().GetKubeconfig(cluster.ClusterId, privilegedusername, privilegedpassword)
		cobra.CheckErr(err)

		yamlbytes := decodeKubeconfig(kubereturn)
		err = kubeconfig.LoadOrNewKubeConfig().MergeYaml(yamlbytes).SetContext(cluster.ClusterName).Write()
		cobra.CheckErr(err)
	}

	_, _ = fmt.Printf("Changed kube context to %s (%s). Happy kubing!\n", colorizeByCluster(&cluster), colorizeByEnv(cluster.ClusterId, cluster.Environment))
	return nil
}

func checkContextAndSwitch(name string) bool {
	kc, err := kubeconfig.LoadFromDefaultFile()
	if err != nil {
		rlog.Error("error loading Kubeconfig", err)
		return false
	}

	isExpired, err := kc.IsExpired(name)
	if err == nil {
		rlog.Error("error", err)
		return false
	}
	if !isExpired {
		err = kc.SetContext(name).Write()
		if err != nil {
			rlog.Error("error writing Kubeconfig", err)
			return false
		}
		return true
	}
	return false
}

func loginToTanzuWorkspace(workspace apicontracts.Workspace) error {
	var privilegedusername, privilegedpassword string

	_, _ = fmt.Printf("\n   %s\n\n", color.HiRedString("Warning: You are logging into a workspace, all changes may affect production"))

	if !checkContextAndSwitch(workspace.Name) {
		var err error
		privilegedusername, privilegedpassword, err = getPrivilegedCreds()
		cobra.CheckErr(err)
		fmt.Println("Logging into workspace " + color.RedString(workspace.Name) + " with user " + color.YellowString(privilegedusername) + "")
		_ = privilegedpassword
		kubereturn, err := config.RorClient.Workspaces().GetKubeconfig(workspace.Name, privilegedusername, privilegedpassword)
		cobra.CheckErr(err)

		yamlbytes := decodeKubeconfig(kubereturn)
		err = kubeconfig.LoadOrNewKubeConfig().MergeYaml(yamlbytes).SetContext(workspace.Name).Write()
		cobra.CheckErr(err)
	}

	_, _ = fmt.Printf("Changed kube context to workspace %s. Happy kubing!\n", color.RedString(workspace.Name))
	return nil
}
