package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"syscall"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"

	rorkubernetesclient "github.com/NorskHelsenett/ror/pkg/clients/kubernetes"
	"github.com/NorskHelsenett/ror/pkg/helpers/k8sportforwarder"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var supportedservices = []string{"argocd", "grafana", "prometheus", "alertmanager", "blackboxexporter"}

func init() {
	rootCmd.AddCommand(portforwardCmd)
}

var portforwardCmd = &cobra.Command{
	Use:     "port-forward service [localport]",
	Aliases: []string{"pf"},
	Short:   "Portforwards a known service to your local machine",
	Long: `
	Portforwards a known service to your local machine

	Supporterd services:
	- argocd
	- grafana
	- prometheus
	- alertmanager
	- blackboxexporter

	Example:

	ror port-forward prometheus

	`,
	Args:              cobra.MaximumNArgs(3),
	ValidArgsFunction: pfCompletion,
	PreRun:            cmdrorclient.SetupRorClient,
	Run:               cmdPortForward,
}

func pfCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return supportedservices, cobra.ShellCompDirectiveNoFileComp
}

func cmdPortForward(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	var id string
	var service string
	var port int32

	service = args[0]

	if len(args) == 2 {
		i, err := strconv.ParseInt(args[1], 10, 32)
		cobra.CheckErr(err)
		port = int32(i)
	}

	if !checkifServiceIsSupported(supportedservices, service) {
		color.Red("Service %s is not supported", service)
		return
	}

	id = viper.GetString(config.LastSessionCluster)

	if id == "" {
		err := fmt.Errorf("missing a clusterId, could not find it in arguments or lastsession")
		cobra.CheckErr(err)
	}

	k8sclientset := rorkubernetesclient.NewK8sClientConfig()
	if k8sclientset == nil {
		panic("failed to initialize kubernetes client")
	}

	forwarder, err := k8sportforwarder.NewPortForwarderFromRorKubernetesClient(k8sclientset)
	cobra.CheckErr(err)

	switch service {
	case "argocd":
		err = forwarder.AddPodByLabels(v1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/name":     "argocd-server",
				"app.kubernetes.io/instance": "argocd",
			},
		}, "argocd")
		forwarder.SetContainerPort(8080)
	case "grafana":
		err = forwarder.AddPodByLabels(v1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/name":     "grafana",
				"app.kubernetes.io/instance": "prometheus",
			},
		}, "prometheus-operator")
		forwarder.SetContainerPort(3000)
	case "prometheus":
		err = forwarder.AddPodByServiceName("prometheus-kube-prometheus-prometheus", "prometheus-operator")
		forwarder.SetContainerPort(9090)
	case "alertmanager":
		err = forwarder.AddPodByServiceName("prometheus-kube-prometheus-alertmanager", "prometheus-operator")

		//err = forwarder.AddPodByName("alertmanager-prometheus-kube-prometheus-alertmanager-0", "prometheus-operator")
		forwarder.SetContainerPort(9093)
	case "blackboxexporter":
		err = forwarder.AddPodByServiceName("prometheus-blackbox-exporter", "prometheus-blackbox-exporter")
		forwarder.SetContainerPort(9115)
	}

	cobra.CheckErr(err)
	if port != 0 {
		forwarder.SetLocalPort(port)
	}

	_, err = config.RorClient.Clusters().GetById(id)
	cobra.CheckErr(err)

	var wg sync.WaitGroup
	wg.Add(1)
	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Bye...")
		close(stopCh)
		wg.Done()
	}()

	go func() {
		err := forwarder.Forward(readyCh, stopCh)
		if err != nil {
			panic(err)
		}
	}()

	<-readyCh

	port, err = forwarder.GetLocalPort()
	cobra.CheckErr(err)
	fmt.Printf("Portforwarding %s to http://localhost:%d\n", service, port)
	url := fmt.Sprintf("http://localhost:%d", port)
	openURL(url)

	wg.Wait()
}

func checkifServiceIsSupported(supported []string, service string) bool {
	return slices.Contains(supported, service)
}
