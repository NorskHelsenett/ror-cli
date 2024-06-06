package outputformatting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"gopkg.in/yaml.v2"
)

func HandleClusterOutputFormatting(outputFormat string, clusters []apicontracts.Cluster) {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		jsoncontent, err := json.MarshalIndent(clusters, "", "\t")
		if err != nil {
			rlog.Fatal("Failed to marshal cluster into jsoncontent: ", err)
		}

		fmt.Println(string(jsoncontent))
	case "yaml":
		yamlContent, err := yaml.Marshal(clusters)
		if err != nil {
			rlog.Fatal("Failed to marshal cluster into yamlContent: ", err)
		}
		fmt.Println(string(yamlContent))
	case "wide":
		clusterWideTabPrinter(clusters)
	case "name":
		for _, cluster := range clusters {
			fmt.Println(cluster.ClusterName)
		}
	case "clusterid":
		for _, cluster := range clusters {
			fmt.Println(cluster.ClusterId)
		}
	case "environment":
		for _, cluster := range clusters {
			fmt.Println(cluster.Environment)
		}
	case "datacenter":
		for _, cluster := range clusters {
			fmt.Println(cluster.Workspace.Datacenter.Name)
		}
	case "created":
		for _, cluster := range clusters {
			fmt.Println(cluster.Created)
		}
	case "lastobserved":
		for _, cluster := range clusters {
			fmt.Println(cluster.LastObserved)
		}
	case "firstobserved":
		for _, cluster := range clusters {
			fmt.Println(cluster.FirstObserved)
		}
	default:
		fmt.Println("unable to match printer for output format. allowed formats: jsoncontent, yamlContent, wide, name, clusterid, environment, datacenter, lastObserved, firstObserved, created")
	}
}

func ClusterTabPrinter(clusters []apicontracts.Cluster) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "ID\tName\tWorkspace\tProvider\tLast observed")
	for _, cluster := range clusters {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v",
			cluster.ClusterId,
			cluster.ClusterName,
			cluster.Workspace.Name,
			cluster.Workspace.Datacenter.Provider,
			cluster.LastObserved)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}

func clusterWideTabPrinter(clusters []apicontracts.Cluster) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "ID\tName\tWorkspace\tEnvironment\tTooling version\tDatacenter\tProvider\tCreated\tLast observed")
	for _, cluster := range clusters {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v",
			cluster.ClusterId,
			cluster.ClusterName,
			cluster.Workspace.Name,
			cluster.Environment,
			cluster.Versions.NhnTooling.Version,
			cluster.Workspace.Datacenter.Name,
			cluster.Workspace.Datacenter.Provider,
			cluster.Created,
			cluster.LastObserved)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}
