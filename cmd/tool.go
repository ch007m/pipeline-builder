package main

import (
	"fmt"
	tool "github.com/ch007m/pipeline-builder"
	"github.com/spf13/cobra"
	"log"
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
		if configurator == "" || output == "" {
			cmd.Usage()
			os.Exit(1)
		}

		// Print the arguments
		fmt.Printf("Configurator: %s\n", configurator)
		fmt.Printf("Output: %s\n", output)

		// Continue to process
		if err := tool.Contribute(configurator, output); err != nil {
			log.Fatal(fmt.Errorf("Unable to generate pipelines ...\n%w", err))
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configurator, "configurator", "c", "", "path of the configurator file")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output where pipelines should be saved")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
