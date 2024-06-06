package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services/cmdrorclient"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"
	"github.com/NorskHelsenett/ror/pkg/rlog"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validWebsevices = []string{"argocd", "grafana"}

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:     "web [cluster] service",
	Aliases: []string{"w"},
	Short:   "Launch the web interface of provided service",
	Long: `
	Display the web interface for the provided service for the loged on cluster.

	Supporterd services:
	- argocd
	- grafana

	Example:

	ror web argocd

	`,
	Args:              cobra.MaximumNArgs(2),
	ValidArgsFunction: webCompletion,
	PreRun:            cmdrorclient.SetupRorClient,
	Run:               cmdWeb,
}

func webCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return validWebsevices, cobra.ShellCompDirectiveNoFileComp
}

func cmdWeb(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	var id string
	var service string
	if len(args) == 2 {
		id = args[0]
		service = args[1]
	} else {
		id = viper.GetString(config.LastSessionCluster)
		service = args[0]
	}

	if !checkifServiceIsSupported(validWebsevices, service) {
		color.Red("Service %s is not supported", service)
		return
	}

	if id == "" {
		rlog.Fatal("error executing command", fmt.Errorf("missing a clusterId, could not find it in arguments or lastsession"))
	}

	cluster, err := config.RorClient.Clusters().GetById(id)
	cobra.CheckErr(err)
	url := getURLfromIngresses(cluster.Ingresses, service)
	fmt.Printf("Opening url: %s\n", url)
	openURL(url)
}

func getURLfromIngresses(ingresses []apicontracts.Ingress, service string) string {
	var ingressname string
	if service == "argocd" {
		ingressname = "argocd-server"
	} else if service == "grafana" {
		ingressname = "grafana-helsenett"
	}

	for _, ingress := range ingresses {
		if ingress.Name == ingressname {
			return fmt.Sprintf("https://%s", ingress.Rules[0].Hostname)
		}
	}
	return ""
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// openURL opens the specified URL in the default browser of the user.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
