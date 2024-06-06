package login

import (
	"os"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RorLogout(cmd *cobra.Command, args []string) {

	all := !cmd.Flags().Changed("ror") && !cmd.Flags().Changed("privileged")

	if cmd.Flags().Changed("privileged") || all {

		color.Yellow("Resetting pam user data")
		viper.Set(config.RorAuthPrivilegedUsername, "")
		viper.Set(config.RorAuthPrivilegedPassword, "")
		viper.Set(config.RorAuthPrivilegedPasswordExpiry, "")
		viper.Set(config.RorAuthCyberarkToken, "")
		viper.Set(config.RorAuthCyberarkExpires, "")
		viper.Set(config.RorAuthCyberarkUsername, "")

	}
	if cmd.Flags().Changed("ror") || all {
		color.Yellow("Resetting ror user data")
		viper.Set(config.RorAuthApiKey, "")
	}

	err := viper.WriteConfig()
	cobra.CheckErr(err)
	os.Exit(0)
}
