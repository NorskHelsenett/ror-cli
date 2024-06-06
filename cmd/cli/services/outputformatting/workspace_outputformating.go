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

func HandleWorkspaceOutputFormatting(outputFormat string, workspaces []apicontracts.Workspace) {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		jsonContent, err := json.MarshalIndent(workspaces, "", "\t")
		if err != nil {
			rlog.Fatal("Failed to marshal workspace(s) into json: ", err)
		}

		fmt.Println(string(jsonContent))
	case "yaml":
		yamlContent, err := yaml.Marshal(workspaces)
		if err != nil {
			rlog.Fatal("Failed to marshal workspace(s) into yamlContent: ", err)
		}
		fmt.Println(string(yamlContent))
	case "wide":
		workspaceWideTabPrinter(workspaces)
	case "name":
		for _, workspace := range workspaces {
			fmt.Println(workspace.Name)
		}
	case "datacenter":
		for _, workspace := range workspaces {
			fmt.Println(workspace.Datacenter.Name)
		}
	case "provider":
		for _, workspace := range workspaces {
			fmt.Println(workspace.Datacenter.Provider)
		}
	case "endpoint":
		for _, workspace := range workspaces {
			fmt.Println(workspace.Datacenter.Provider)
		}
	default:
		fmt.Println("unable to match printer for output format. allowed formats: json, yamlContent, wide, name, datacenter, provider, endpoint")
	}
}

func WorkspaceTabPrinter(workspaces []apicontracts.Workspace) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tDatacenter")
	for _, workspace := range workspaces {
		data := fmt.Sprintf("%v\t%v",
			workspace.Name,
			workspace.Datacenter.Name)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}

func workspaceWideTabPrinter(workspaces []apicontracts.Workspace) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tDatacenter\tProvider\tEndpoint")
	for _, workspace := range workspaces {
		data := fmt.Sprintf("%v\t%v\t%v\t%v",
			workspace.Name,
			workspace.Datacenter.Name,
			workspace.Datacenter.Provider,
			workspace.Datacenter.APIEndpoint)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}
