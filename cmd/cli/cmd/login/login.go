package login

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"

	"github.com/NorskHelsenett/ror/pkg/clients/cyberark"

	"github.com/NorskHelsenett/ror/pkg/rlog"
	"github.com/fatih/color"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/NorskHelsenett/ror/pkg/apicontracts"
)

func RorLogin(cmd *cobra.Command, args []string) {

	var sorted []interface{}

	isWorkspace, _ := cmd.Flags().GetBool("workspace")
	// In case the user specifies what cluster/workspace they want to use manually. ror login {clusterId}
	if len(args) == 1 {
		id := args[0]

		if isWorkspace {
			workspace, err := config.RorClient.Workspaces().GetByName(id)
			cobra.CheckErr(err)

			if workspace.Name == "" {
				_, _ = fmt.Fprintf(os.Stderr, "invalid workspace %s\n", id)
				return
			}
			LogIntoWorkspace(*workspace)
			return
		}

		cluster, err := config.RorClient.Clusters().GetById(id)

		if err != nil {
			cobra.CheckErr(err)
		}

		if cluster.ClusterId == "" {
			_, _ = fmt.Fprintf(os.Stderr, "invalid cluster %s\n", id)
			return
		}

		logIntoCluster(*cluster)
		return
	}

	var templates promptui.SelectTemplates
	if !isWorkspace {
		sorted = getClustersSorted()
		templates = promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ .ClusterName | cyan }} ({{ .ClusterId }})",
			Selected: "{{ .ClusterName | red}}",
			Inactive: "{{ .ClusterName }}",
		}
		promptSearchAndLogin(sorted, &templates, isWorkspace)
		return
	}

	sorted = getWorkspacesSorted()
	templates = promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "{{ .Name | cyan }}",
		Selected: "{{ .Name | red | cyan }}",
		Inactive: "{{ .Name }}",
	}
	promptSearchAndLogin(sorted, &templates, isWorkspace)
}

func logIntoCluster(cluster apicontracts.Cluster) {

	switch cluster.Workspace.Datacenter.Provider {
	case "tanzu":
		err := loginToTanzuCluster(cluster)
		cobra.CheckErr(err)
		viper.Set(config.LastSessionCluster, cluster.ClusterId)
		viper.Set(config.LastSessionWorkspace, cluster.Workspace.Name)
		err = viper.WriteConfig()
		cobra.CheckErr(err)
		return
	default:
		fmt.Println("Login does not support provider " + cluster.Workspace.Datacenter.Provider + " yet")
		return
	}
}

func colorizeByCluster(cluster *apicontracts.Cluster) string {
	return colorizeByEnv(cluster.ClusterName, cluster.Environment)
}

func colorizeByEnv(text string, env string) string {
	if env == "prod" || env == "mgmt" {
		return color.RedString(text)
	}
	if env == "dev" {
		return color.BlueString(text)
	}
	if env == "test" {
		return color.GreenString(text)
	}
	if env == "qa" {
		return color.YellowString(text)
	}
	if env == "kurs" {
		return color.MagentaString(text)
	}
	if env == "mgmt" {
		return color.RedString(text)
	}
	return text
}

func getPrivilegedCreds() (string, string, error) {
	var privilegedusername, privilegedpassword string
	var err error

	// If user is already authenticated and we have a valid password we can use that
	if checkPrivilegedUserStilValid() && viper.GetBool(config.RorAuthClientConfigSecure) {
		privilegedusername = viper.GetString(config.RorAuthPrivilegedUsername)
		privilegedpassword = viper.GetString(config.RorAuthPrivilegedPassword)
		return privilegedusername, privilegedpassword, nil
	}

	cyberarkcli, err := cyberark.NewCyberarkClient("https://tilgang.pam.nhn.no", config.CyberarkValidDomains...)
	cobra.CheckErr(err)

	// Handle the case where we dont have access to the pam server

	if !cyberarkcli.Ping() {

		if checkPrivilegedUserStilValid() && viper.GetBool(config.RorAuthClientConfigSecure) {
			privilegedusername = viper.GetString(config.RorAuthPrivilegedUsername)
			privilegedpassword = viper.GetString(config.RorAuthPrivilegedPassword)
		} else {

			color.Yellow("Warning: Could not connect to pam portal, you have to get the password from the portal your self.")
			privilegedusername = getPrivilegedUsername()
			privilegedpassword = PromtPassword("Enter password from the pam portal for user " + privilegedusername + "")
		}

		if viper.GetBool(config.RorAuthClientConfigSecure) {
			viper.Set(config.RorAuthPrivilegedPassword, privilegedpassword)
			viper.Set(config.RorAuthPrivilegedPasswordExpiry, twoAm())
		}

		err = viper.WriteConfig()
		cobra.CheckErr(err)

		return privilegedusername, privilegedpassword, nil
	}

	// Get username from config or prompt, wil persist an reuse
	pamusername := getPamUsername()

	var secrets *[]cyberark.CyberarkSecret

	// Check if we have a valid token
	if !(viper.IsSet(config.RorAuthCyberarkToken) && viper.IsSet(config.RorAuthCyberarkExpires) && viper.GetTime(config.RorAuthCyberarkExpires).After(time.Now())) {

		pampassword := PromtPassword("Enter pam password")
		color.Yellow("Please complete the 2fa challenge")
		token, expires, err := cyberarkcli.Authenticate(pamusername, pampassword)
		if err != nil {
			color.Red("Could not authenticate user " + pamusername + " to pam portal, or complete 2fa challenge\n Please try again...")
			//cobra.CheckErr(err)
		}

		viper.Set(config.RorAuthCyberarkToken, token)
		viper.Set(config.RorAuthCyberarkExpires, expires)
		err = viper.WriteConfig()
		if err != nil {
			cobra.CheckErr(err)
		}
	} else {
		cyberarkcli.SetToken(viper.GetString(config.RorAuthCyberarkToken))
	}

	if viper.IsSet(config.RorAuthPrivilegedUsername) && viper.GetString(config.RorAuthPrivilegedUsername) != "" {
		secret, err := cyberarkcli.GetSecret(viper.GetString(config.RorAuthPrivilegedUsername))
		cobra.CheckErr(err)
		secrets = &[]cyberark.CyberarkSecret{*secret}
	} else {
		secrets, err = cyberarkcli.GetSecrets()
		cobra.CheckErr(err)
	}

	if len(*secrets) == 0 {
		fmt.Println("No secrets found for user " + pamusername + " in pam portal")
		os.Exit(1)
	}
	if len(*secrets) == 1 {
		secret := (*secrets)[0]
		privilegedusername = fmt.Sprintf("%s@%s", secret.UserName, secret.Address)
		privilegedpassword, err = cyberarkcli.GetPassword((*secrets)[0].ID)
		cobra.CheckErr(err)

		if viper.GetBool(config.RorAuthClientConfigSecure) {
			viper.Set(config.RorAuthPrivilegedPassword, privilegedpassword)
			viper.Set(config.RorAuthPrivilegedPasswordExpiry, twoAm())
		}
		err = viper.WriteConfig()
		cobra.CheckErr(err)
		return privilegedusername, privilegedpassword, err
	}
	elements := []cyberark.CyberarkSecret{}
	for _, secret := range *secrets {
		secret.Displayname = fmt.Sprintf("%s@%s", secret.UserName, secret.Address)

		elements = append(elements, secret)
	}
	userprompt := promptui.Select{
		Label: "Select privileged user",
		Items: elements,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "{{ .Displayname | cyan }}",
			Selected: "{{ .Displayname | red | cyan }}",
			Inactive: "{{ .Displayname }}",
		},
	}
	i, _, err := userprompt.Run()
	cobra.CheckErr(err)
	privilegedusername = elements[i].Displayname

	viper.Set(config.RorAuthPrivilegedUsername, privilegedusername)

	privilegedpassword, err = cyberarkcli.GetPassword(elements[i].ID)
	cobra.CheckErr(err)
	if viper.GetBool(config.RorAuthClientConfigSecure) {
		viper.Set(config.RorAuthPrivilegedPassword, privilegedpassword)
		viper.Set(config.RorAuthPrivilegedPasswordExpiry, twoAm())
	}
	err = viper.WriteConfig()
	cobra.CheckErr(err)

	return privilegedusername, privilegedpassword, err
}

// we dont have to login if the cluster token is still valid

func checkPrivilegedUserStilValid() bool {
	if viper.IsSet(config.RorAuthPrivilegedUsername) && viper.IsSet(config.RorAuthPrivilegedPassword) && viper.IsSet(config.RorAuthPrivilegedPasswordExpiry) {
		if viper.GetTime(config.RorAuthPrivilegedPasswordExpiry).After(time.Now()) {
			return true
		}
	}

	return false
}
func twoAm() time.Time {
	now := time.Now()
	var tomorrow time.Time
	if now.Hour() >= 0 && now.Hour() <= 2 {
		tomorrow = time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, time.Now().Location())
	} else {
		tomorrow = time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 1)
	}
	return tomorrow
}
func decodeKubeconfig(kubeconfig *apicontracts.ClusterKubeconfig) []byte {
	sDec, err := base64.StdEncoding.DecodeString(kubeconfig.Data)
	cobra.CheckErr(err)
	return sDec

}

func getPamUsername() string {
	var username string

	username = viper.GetString(config.RorAuthCyberarkUsername)
	if username != "" {
		return username
	}
	fmt.Println()
	user, err := user.Current()
	cobra.CheckErr(err)

	username = user.Username

	if strings.Contains(username, "\\") {
		parts := strings.Split(username, "\\")
		username = parts[1]
	}

	if strings.Contains(username, "@") {
		parts := strings.Split(username, "@")
		username = parts[0]
	}

	username = strings.TrimPrefix(username, "t1-")

	prompt := promptui.Prompt{
		Label: "Enter pam username",
		Validate: func(input string) error {
			regex := regexp.MustCompile(`([a-z0-9]+)`)
			if !regex.MatchString(input) {
				return errors.New("wrong username format")
			}
			return nil
		},
		Default: username,
	}
	username, err = prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			cobra.CheckErr(err)
		}

	}
	viper.Set(config.RorAuthCyberarkUsername, username)
	_ = viper.WriteConfig()
	return username
}

func getPrivilegedUsername() string {
	var username string

	username = viper.GetString(config.RorAuthPrivilegedUsername)
	if username != "" {
		return username
	}
	fmt.Println()
	user, err := user.Current()
	cobra.CheckErr(err)

	username = user.Username
	if !strings.HasPrefix(username, "t1-") {
		username = "t1-" + username
	}
	// add select
	username = username + "@cloud.nhn.no"

	prompt := promptui.Prompt{
		Label: "Enter privileged username (t1-username@domain.nhn.no)",
		Validate: func(input string) error {
			regex := regexp.MustCompile(`t1-([a-z0-9]+)@([a-z]+).nhn.no`)
			if !regex.MatchString(input) {
				return errors.New("wrong username format")
			}
			return nil
		},
		Default: username,
	}
	username, err = prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			cobra.CheckErr(err)
		}

	}
	viper.Set(config.RorAuthPrivilegedUsername, username)
	_ = viper.WriteConfig()
	return username
}

func PromtPassword(descr string) string {

	prompt := promptui.Prompt{
		Label:       descr,
		Mask:        '*',
		HideEntered: true,
	}

	pampassword, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			cobra.CheckErr(err)
		}
	}
	return pampassword
}
func LogIntoWorkspace(workspace apicontracts.Workspace) {
	err := loginToTanzuWorkspace(workspace)
	cobra.CheckErr(err)
}

func promptSearchAndLogin(elements []interface{}, templates *promptui.SelectTemplates, isWorkspace bool) {
	label := "Cluster"
	if isWorkspace {
		label = "Workspace"
	}
	searcher := func(input string, index int) bool {
		elem := elements[index]

		var id string
		if isWorkspace {
			id = (elem.(apicontracts.Workspace)).Name
		} else {
			id = (elem.(apicontracts.Cluster)).ClusterName
		}

		name := strings.Replace(strings.ToLower(id), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	vimMode := false

	if viper.GetBool(config.Vim) {
		vimMode = true
	}

	prompt := promptui.Select{
		Label:             label,
		Items:             elements,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		HideSelected:      true,
		StartInSearchMode: true,
		IsVimMode:         vimMode,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return
		}
		_, _ = fmt.Fprintf(os.Stderr, "could not display cluster search prompt")
		rlog.Fatal("Prompt failed", err)
	}

	if isWorkspace {
		LogIntoWorkspace(elements[i].(apicontracts.Workspace))
	} else {
		logIntoCluster(elements[i].(apicontracts.Cluster))
	}
}

func getClustersSorted() []interface{} {
	var sorted []interface{}
	clusters, err := config.RorClient.Clusters().GetAll()
	if err != nil {
		cobra.CheckErr(err)
	}
	if last := viper.GetString(config.LastSessionCluster); last != "" {
		elem := apicontracts.Cluster{}

		for _, c := range *clusters {
			if c.ClusterId == last {
				elem = c
			}
		}

		sorted = append(sorted, elem)

		for _, c := range *clusters {
			if c.ClusterName != elem.ClusterName {
				sorted = append(sorted, c)
			}
		}

	} else {
		for _, c := range *clusters {
			sorted = append(sorted, c)
		}
	}
	return sorted
}

func getWorkspacesSorted() []interface{} {
	var sorted []interface{}
	workspaces, err := config.RorClient.Workspaces().GetAll()
	if workspaces == nil {
		return sorted
	}
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "could not get workspaces from ror")
	}
	if last := viper.GetString(config.LastSessionWorkspace); last != "" {
		elem := apicontracts.Workspace{}

		for _, w := range *workspaces {
			if w.Name == last {
				elem = w
			}
		}

		sorted = append(sorted, elem)
		for _, w := range *workspaces {
			if w.Name != elem.Name {
				sorted = append(sorted, w)
			}
		}
	} else {
		for _, w := range *workspaces {
			sorted = append(sorted, w)
		}
	}

	return sorted
}
