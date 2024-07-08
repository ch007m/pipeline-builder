package main

import (
	tool "github.com/ch007m/pipeline-builder/builder"
	"github.com/ch007m/pipeline-builder/logging"
	"github.com/spf13/cobra"
	"os"
)

var configurator string
var output string

var rootCmd = &cobra.Command{
	Use:   "pipeline-builder",
	Short: "A tekton pipeline builder",
	Long:  `A tekton pipeline builder able to create from templates and a configurator files pipelines ans tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if required flags are provided
		if configurator == "" {
			logging.Logger.Warn("The argument: -c <config_file> is missing.")
			os.Exit(1)
		}

		if output == "" {
			logging.Logger.Warn("The argument: -o <output_dir> is missing.")
			os.Exit(1)
		}

		// Print the arguments
		logging.Logger.Debug("Configurator: %s\n", configurator)
		logging.Logger.Debug("Output: %s\n", output)

		// Continue to process
		if err := tool.Contribute(configurator, output); err != nil {
			logging.Logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configurator, "configurator", "c", "", "path of the configurator file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output where pipelines should be saved")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}
}
