package cmd

import (
	"fmt"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	promptCmd.Flags().BoolVar(&isPS1Output, "PS1", false, "Format output for PS1, i.e. wrap colour bytes with '\\[' & '\\]'")
	rootCmd.AddCommand(promptCmd)
}

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prints cluster info for prompts",
	Long: `
	Prints out information about cluster context for users to use in their bash,zsh,etc. prompts.

	Here is an example of a very basic prompt utilizing this command:

		PS1='$(ror prompt --PS1 2>/dev/null)\n\u@\h \W \$ '

	Paste this snippet into your .bashrc or .zshrc file and source it and you are good to go. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		printCubeContext(isPS1Output)
	},
}

func printCubeContext(isPS1Output bool) {
	context := viper.GetString(config.LastSessionCluster)
	fmt.Print(ColorStringByLastEnv(context, isPS1Output))
}

// colors the input string based on the lastsession environment set in the
// ror-cli configfile. Last env is also the current env.
// Can also format for PS1 (command prompt)
func ColorStringByLastEnv(name string, isPS1Output bool) string {
	environment := viper.GetString(config.LastSessionEnvironment)
	const ansiReset = "\u001b[0m"
	var ansiColor string

	switch environment {

	case "prod":
		ansiColor = "\u001b[31m" //red
	case "mgmt":
		ansiColor = "\u001b[31m" //red
	case "qa":
		ansiColor = "\u001b[33m" //yellow
	case "test":
		ansiColor = "\u001b[32m" //green
	case "lab":
		ansiColor = "\u001b[34m" //blue
	case "dev":
		ansiColor = "\u001b[34m" //blue
	default:
		ansiColor = "\u001b[37m" //white
	}

	var outputFormat string
	switch isPS1Output {
	case true:
		outputFormat = "\\[%s\\]%s\\[%s\\]" // double backslash to escape such that final output is \[%s\]%s\[%s\], this escapes byte counting colours in bash prompt
	case false:
		outputFormat = "%s%s%s"
	}

	return fmt.Sprintf(outputFormat, ansiColor, name, ansiReset)
}
