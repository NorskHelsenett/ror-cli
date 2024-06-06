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

func HandleDatacenterOutputFormatting(outputFormat string, datacenters *[]apicontracts.Datacenter) {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		content, err := json.MarshalIndent(datacenters, "", "\t")
		if err != nil {
			rlog.Fatal("Failed to marshal datacenter(s) into content: ", err)
		}

		fmt.Println(string(content))
	case "yaml":
		yamlContent, err := yaml.Marshal(datacenters)
		if err != nil {
			rlog.Fatal("Failed to marshal datacenter(s) into yamlContent: ", err)
		}
		fmt.Println(string(yamlContent))
	case "wide":
		datacenterWideTabPrinter(*datacenters)
	case "name":
		for _, datacenter := range *datacenters {
			fmt.Println(datacenter.Name)
		}
	case "provider":
		for _, datacenter := range *datacenters {
			fmt.Println(datacenter.Provider)
		}
	case "region":
		for _, datacenter := range *datacenters {
			fmt.Println(datacenter.Location.Region)
		}
	case "country":
		for _, datacenter := range *datacenters {
			fmt.Println(datacenter.Location.Country)
		}
	case "endpoint":
		for _, datacenter := range *datacenters {
			fmt.Println(datacenter.APIEndpoint)
		}
	default:
		fmt.Println("unable to match printer for output format. allowed formats: content, yamlContent, wide, name, provider, region, country, endpoint")
	}
}

func DatacenterTabPrinter(datacenters []apicontracts.Datacenter) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tProvider")
	for _, datacenter := range datacenters {
		data := fmt.Sprintf("%v\t%v",
			datacenter.Name,
			datacenter.Provider)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}

func datacenterWideTabPrinter(datacenters []apicontracts.Datacenter) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tProvider\tRegion\tCountry\tEnpoint")
	for _, datacenter := range datacenters {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v",
			datacenter.Name,
			datacenter.Provider,
			datacenter.Location.Region,
			datacenter.Location.Country,
			datacenter.APIEndpoint)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}
