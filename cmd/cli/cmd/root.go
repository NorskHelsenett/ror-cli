package cmd

import (
	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile          string
	accesstoken      string
	verbose          bool
	suppressPrompts  bool
	debug            bool
	isPS1Output      bool
	outputFormat     string
	paginationLimit  int
	paginationOffset int
	rootCmd          = &cobra.Command{
		Use:   "ror ",
		Short: "Get status related to this NHN-ROR-CLI",
		Long:  `NHN ROR CLI `,
		//Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
				return nil
			}
			return nil
		},
	}
)

// Execute executes the root command .
func Execute() error {
	return rootCmd.Execute()
}

// the init function
func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(rlog.InitializeRlog)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ror/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&accesstoken, "token", "", "Access token")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&suppressPrompts, "silent", "s", false, "no output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	_ = viper.BindPFlag("accessToken", rootCmd.PersistentFlags().Lookup("accessToken"))
}

func initConfig() {
	config.Load(cfgFile)
}
