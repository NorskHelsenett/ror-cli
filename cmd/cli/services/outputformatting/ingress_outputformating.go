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

type IngressEntry struct {
	Name      string
	Namespace string
	Host      string
	Service   string
	Type      string
	Status    string
	Ip        string
}

func HandleIngressOutputFormatting(outputFormat string, ingresses []IngressEntry, rawIngresses []apicontracts.Ingress) {
	outputFormat = strings.ToLower(outputFormat)
	switch outputFormat {
	case "json":
		jsonContent, err := json.MarshalIndent(rawIngresses, "", "\t")
		if err != nil {
			rlog.Fatal("Failed to marshal ingress(s) into jsonContent: ", err)
		}

		fmt.Println(string(jsonContent))
	case "yaml":
		yamlContent, err := yaml.Marshal(rawIngresses)
		if err != nil {
			rlog.Fatal("Failed to marshal ingress(s) into yamlContent: ", err)
		}
		fmt.Println(string(yamlContent))
	case "wide":
		ingressWideTabPrinter(ingresses)
	case "name":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Name)
		}
	case "namespace":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Namespace)
		}
	case "host":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Host)
		}
	case "service":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Service)
		}
	case "type":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Type)
		}
	case "status":
		for _, ingress := range ingresses {
			fmt.Println(ingress.Status)
		}
	default:
		fmt.Println("unable to match printer for output format. allowed formats: jsonContent, yamlContent, wide, name, namespace, host, service, type, status")
	}
}

func IngressTabPrinter(ingresses []IngressEntry) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tNamespace\tHost\tStatus\tIp")
	for _, ingress := range ingresses {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v",
			ingress.Name,
			ingress.Namespace,
			ingress.Host,
			ingress.Status,
			ingress.Ip)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}

func ingressWideTabPrinter(ingresses []IngressEntry) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, "Name\tNamespace\tHost\tStatus\tService\tType")
	for _, ingress := range ingresses {
		data := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v",
			ingress.Name,
			ingress.Namespace,
			ingress.Host,
			ingress.Status,
			ingress.Service,
			ingress.Type)

		_, _ = fmt.Fprintln(writer, data)
	}

	_ = writer.Flush()
}
