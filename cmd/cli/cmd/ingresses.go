package cmd

import (
	"fmt"
	"strings"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/outputformatting"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ingressesCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Decide how you want data to be displayed (json/default)")
	rootCmd.AddCommand(ingressesCmd)
}

var ingressesCmd = &cobra.Command{
	Use:     "ingresses",
	Aliases: []string{"i"},
	Short:   "Get details about all ingresses",
	Long: `
	Display ingresses of a cluster you have access too.
	The cluster is inferred either by the last cluster you was logged into or 
	a command line argument

	Example:

	ror ingresses tooling-staging.trd1-nhn-tooling

	`,
	PreRun: cmdrorclient.SetupRorClient,
	Run:    cmdIngresses,
}

func cmdIngresses(cmd *cobra.Command, args []string) {
	var id string

	if len(args) > 0 {
		id = args[0]
	} else {
		id = viper.GetString(config.LastSessionCluster)
	}

	if id == "" {
		rlog.Fatal("error executing command", fmt.Errorf("missing a clusterId, could not find it in arguments or lastsession"))
	}

	cluster, err := config.RorClient.Clusters().GetById(id)
	cobra.CheckErr(err)

	ingresses := cluster.Ingresses

	var entries []outputformatting.IngressEntry
	for _, ingress := range ingresses {
		for _, rule := range ingress.Rules {
			for _, path := range rule.Paths {
				var ingressType string

				if ingress.Class == "avi-ingress-class-helsenett" {
					ingressType = "helsenett"
				} else if ingress.Class == "avi-ingress-class-datacenter" {
					ingressType = "datacenter"
				} else if ingress.Class == "avi-ingress-class-internett" {
					// Assume the ingress is exposed to the internet if nothing else is provided
					ingressType = "internett"
				} else {
					ingressType = ingress.Class
				}

				var status string

				if ingress.Health == apicontracts.HealthHealthy {
					status = "healthy"
				} else if ingress.Health == apicontracts.HealthUnhealthy {
					status = "unhealthy"
				} else if ingress.Health == apicontracts.HealthBad {
					status = "bad"
				} else {
					status = "unknown"
				}

				ip := strings.Join(rule.IPAddresses, ", ")

				entry := outputformatting.IngressEntry{
					Name:      ingress.Name,
					Namespace: ingress.Namespace,
					Host:      rule.Hostname,
					Service:   path.Service.Name,
					Type:      ingressType,
					Status:    status,
					Ip:        ip,
				}
				entries = append(entries, entry)
			}
		}
	}

	if !cmd.Flags().Changed("output") {
		outputformatting.IngressTabPrinter(entries)
		return
	}

	outputformatting.HandleIngressOutputFormatting(outputFormat, entries, ingresses)
}
