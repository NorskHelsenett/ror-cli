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

func HandleProjectOutputFormatting(outputFormat string, projects []apicontracts.Project) {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		jsonContent, err := json.MarshalIndent(projects, "", "\t")
		if err != nil {
			rlog.Fatal("Failed to marshal project into jsonContent: ", err)
		}

		fmt.Println(string(jsonContent))
	case "yaml":
		yamlContent, err := yaml.Marshal(projects)
		if err != nil {
			rlog.Fatal("Failed to marshal project into yamlContent: ", err)
		}
		fmt.Println(string(yamlContent))
	case "wide":
		projectWideTabPrinter(projects)
	case "name":
		for _, project := range projects {
			fmt.Println(project.Name)
		}
	default:
		fmt.Println("unable to match printer for output format. allowed formats: jsonContent, yamlContent, wide, name")
	}
}

func ProjectTabPrinter(projects []apicontracts.Project) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tActive")
	for _, project := range projects {
		data := fmt.Sprintf("%v\t%v",
			project.Name,
			project.Active)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}

func projectWideTabPrinter(projects []apicontracts.Project) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tActive\tBilling\tRoles\tCreated")
	for _, project := range projects {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v",
			project.Name,
			project.Active,
			project.ProjectMetadata.Billing,
			project.ProjectMetadata.Roles,
			project.Created)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}
