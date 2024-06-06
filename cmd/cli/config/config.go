package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/models"

	"github.com/NorskHelsenett/ror/pkg/clients/rorclient"

	"github.com/NorskHelsenett/ror/pkg/config/rorversion"

	"github.com/NorskHelsenett/ror/pkg/apicontracts/v2/apicontractsv2self"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viper key consts
const ApiRor string = "apiconfig.ror"
const ApiDex string = "apiconfig.dex"
const ApiVsphere string = "apiconfig.vsphere"

const LogLevel string = "log_level"
const LogOutput string = "log_output"
const LogOutputError string = "log_output_error"

const LastSessionCluster string = "lastsession.cluster"
const LastSessionWorkspace string = "lastsession.workspace"
const LastSessionEnvironment string = "lastsession.environment"

const RorAuthApiKey string = "rorauth.apikey"
const RorAuthCyberarkUsername string = "rorauth.cyberark.username"
const RorAuthCyberarkToken string = "rorauth.cyberark.token"
const RorAuthCyberarkExpires string = "rorauth.cyberark.expires"
const RorAuthPrivilegedUsername string = "rorauth.privileged.username"
const RorAuthPrivilegedPassword string = "rorauth.privileged.password"
const RorAuthPrivilegedPasswordExpiry string = "rorauth.privileged.expires"
const RorAuthClientConfigSecure string = "rorauth.client.secure"
const AuthId string = "auth.id"
const AuthType string = "auth.type"
const AuthToken string = "auth.token"
const AuthExpiry string = "auth.expiry"
const AuthRefresh string = "auth.refresh"

const VsphereExpiry string = "vsphere.expiry"
const VspherePassword string = "vsphere.password"

const ServerVersion string = "server.version"

const Vim string = "vim"

var Config models.ApiConfig

var validate *validator.Validate

var Version = "0.0.0"
var Commit = "ffffff"

var Role = "ror-cli"
var Authinfo apicontractsv2self.SelfData

var RorClient *rorclient.RorClient

var RorVersion rorversion.RorVersion

// ror-config defaults
var defaults = models.CliConfig{
	Log_level:        "info",
	Log_output:       getDefaultLogfilePath(),
	Log_output_error: "",
	Apiconfig:        defaultAPIs,
}

var defaultAPIs = models.ApiConfig{
	Ror:     "https://api.ror.sky.test.nhn.no",
	Dex:     "https://auth.sky.nhn.no/dex",
	Vsphere: "https://ptr1-w02-cl01-api.sdi.nhn.no",
}
var CyberarkValidDomains = []string{"cloud.nhn.no", "drift.nhn.no"}

func Load(cfgFile string) {
	configPath := getDefaultConfigPath()

	if cfgFile != "" {
		configPath = cfgFile
		_, err := os.Stat(configPath)
		if errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(os.Stderr, "failed to read the config file in path %s: %s\n", configPath, err)
			os.Exit(1)
		}
	}

	viper.SetConfigFile(configPath)

	info, err := os.Stat(configPath)
	if errors.Is(err, os.ErrNotExist) {
		CreateDefaultConfigFile()
		return
	}

	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read the config file in path %s\n", configPath)
		os.Exit(1)
	}

	validateConfig(configPath)

	if info.Mode() != 0600 && runtime.GOOS != "windows" {
		_, _ = fmt.Fprintln(os.Stderr, "config file does not have strict enough permissions, wont allow persisting of privileged credentials")
		viper.Set(RorAuthClientConfigSecure, false)
	} else {
		viper.Set(RorAuthClientConfigSecure, true)
	}

	RorVersion = rorversion.NewRorVersion(Version, Commit)
	viper.SetDefault(ServerVersion, "notconnected")
	viper.AutomaticEnv()
}

func validateConfig(cfgFile string) {
	validate := validator.New()

	var config models.CliConfig

	err := viper.Unmarshal(&config)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "could not validate config")
		os.Exit(1)
	}

	err = validate.Struct(config)
	if err != nil {
		var keys []string
		for _, err := range err.(validator.ValidationErrors) {
			keys = append(keys, err.Field())
		}

		err = addMissingRequiredKeys(keys)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "could not add missing required fields to config")
			os.Exit(1)
		}
		return
	}
}

// CreateDefaultConfigFile creates a new, or overwrites an old config file
func CreateDefaultConfigFile() {
	configDirPath := getDefaultConfigDirPath()
	configFullPath := getDefaultConfigPath()

	_ = os.MkdirAll(configDirPath, os.ModePerm)

	config := defaults

	configBytes, err := json.Marshal(config)
	cobra.CheckErr(err)

	_ = viper.ReadConfig(bytes.NewBuffer(configBytes))

	err = viper.WriteConfigAs(configFullPath)
	cobra.CheckErr(err)

	err = os.Chmod(configFullPath, 0600)
	cobra.CheckErr(err)
}

// takes a set of keys and adds their default value to the config
// this is a softer approach than recreating the config entierly
// this will only set required keys
func addMissingRequiredKeys(keys []string) error {
	for _, key := range keys {

		switch key {
		case "Log_level":
			viper.Set("log_level", defaults.Log_level)
		case "Log_output":
			viper.Set("log_output", defaults.Log_output)
		case "Ror":
			viper.Set("apiconfig.ror", defaultAPIs.Ror)
		case "Dex":
			viper.Set("apiconfig.dex", defaultAPIs.Dex)
		case "Vsphere":
			viper.Set("apiconfig.vsphere", defaultAPIs.Vsphere)
		case "Vim":
			{
				viper.Set("vim", true)
			}
		}
	}
	return viper.WriteConfig()
}

func getDefaultConfigDirPath() string {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return path.Join(homeDir, ".ror")
}

func getDefaultConfigPath() string {
	configDir := getDefaultConfigDirPath()
	return path.Join(configDir, "config.yaml")
}

func getDefaultLogfilePath() string {
	configDir := getDefaultConfigDirPath()
	return path.Join(configDir, "log")
}
