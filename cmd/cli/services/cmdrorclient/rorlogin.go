package cmdrorclient

import (
	"fmt"
	"os"
	"time"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/responses"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/services"

	"github.com/NorskHelsenett/ror/pkg/clients/rorclient"
	"github.com/NorskHelsenett/ror/pkg/clients/rorclient/transports/resttransport"
	"github.com/NorskHelsenett/ror/pkg/clients/rorclient/transports/resttransport/httpauthprovider"
	"github.com/NorskHelsenett/ror/pkg/clients/rorclient/transports/resttransport/httpclient"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SetupRorNonAuthClient(cmd *cobra.Command, args []string) {
	transport := resttransport.NewRorHttpTransport(&httpclient.HttpTransportClientConfig{
		BaseURL:      viper.GetString(config.ApiRor),
		AuthProvider: httpauthprovider.NewNoAuthprovider(),
		Version:      config.RorVersion,
		Role:         config.Role,
	})
	config.RorClient = rorclient.NewRorClient(transport)

	serverversion, err := config.RorClient.Info().GetVersion()
	if err != nil {
		viper.Set(config.ServerVersion, "Unknown")
	}
	viper.Set(config.ServerVersion, serverversion)
}

func SetupRorClient(cmd *cobra.Command, args []string) {
	var authenticatedok bool
	var err error

	if viper.IsSet(config.RorAuthApiKey) {
		transport := resttransport.NewRorHttpTransport(&httpclient.HttpTransportClientConfig{
			BaseURL:      viper.GetString(config.ApiRor),
			AuthProvider: httpauthprovider.NewAuthProvider(httpauthprovider.AuthPoviderTypeAPIKey, viper.GetString(config.RorAuthApiKey)),
			Version:      config.RorVersion,
			Role:         config.Role,
		})
		config.RorClient = rorclient.NewRorClient(transport)

		err := config.RorClient.Ping()
		if err != nil {
			err = fmt.Errorf("could not connect to ror api: %v", err)
			cobra.CheckErr(err)
		}

		config.Authinfo, err = config.RorClient.Self().Get()
		if err != nil {
			//fmt.Println(err.Error())
			authenticatedok = false
		} else {
			authenticatedok = true
		}

	}
	if !authenticatedok {
		tokens, err := authenticateRor()
		cobra.CheckErr(err)
		transport := resttransport.NewRorHttpTransport(&httpclient.HttpTransportClientConfig{
			BaseURL:      viper.GetString(config.ApiRor),
			AuthProvider: httpauthprovider.NewAuthProvider(httpauthprovider.AuthProviderTypeBearer, tokens.AccessToken),
			Version:      config.RorVersion,
			Role:         config.Role,
		})
		config.RorClient = rorclient.NewRorClient(transport)
		config.Authinfo, err = config.RorClient.Self().Get()
		if err != nil {
			fmt.Println("Could not authenticate with oicd")
			cobra.CheckErr(err)
		}

		// Todo move to function
		hostname, _ := os.Hostname()
		keyname := fmt.Sprintf("ror-cli: %s", hostname)

		apikey, err := config.RorClient.Self().CreateOrUpdateApiKey(keyname, 7776000)
		if err != nil {
			cobra.CheckErr(err)
		}
		viper.Set(config.RorAuthApiKey, apikey)
		err = viper.WriteConfig()
		if err != nil {
			cobra.CheckErr(err)
		}
	}

	serverversion, err := config.RorClient.Info().GetVersion()
	if err != nil {
		err2 := fmt.Errorf("Could not connect to ror api: %v", err)
		cobra.CheckErr(err2)
	}
	viper.Set(config.ServerVersion, serverversion)
	if cmd.Flag("verbose").Value.String() == "true" {
		_, _ = fmt.Printf("Connected to ROR API: %s, Client: %s Server: %s\n", viper.GetString(config.ApiRor), config.RorVersion.Version, viper.GetString(config.ServerVersion))
	}
}

func authenticateRor() (*responses.TokenResponse, error) {
	writer := os.Stdout
	code, err := services.GetDeviceCode()
	if err != nil {
		return nil, err
	}
	fmt.Println("To authenticate open this uri in your nearest browser: ")
	_, _ = fmt.Printf("Url: %v \n", code.VerificationUriComplete)

	_, _ = fmt.Fprint(writer, "Waiting for you to authenticate")

	stop := make(chan bool)
	go startSpinner(stop)

	tokens := services.GetJWToken(code)

	stop <- true
	close(stop)

	return tokens, nil
}
func startSpinner(stop chan bool) {
	for {
		select {
		case <-stop:
			fmt.Print("\n")
			return
		default:
			_, _ = fmt.Printf(".")
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
