package main

import (
	"fmt"
	tool "github.com/ch007m/pipeline-builder"
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	flagSet := pflag.NewFlagSet("octo", pflag.ExitOnError)
	configurator := flagSet.String("configurator", "", "path to input configurator")
	output := flagSet.String("output", "", "path to output directory")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatal(fmt.Errorf("unable to parse flags\n%w", err))
	}

	if configurator == nil {
		log.Fatal("--configurator is required")
	}

	if output == nil {
		log.Fatal("--output is required")
	}

	if err := tool.Contribute(*configurator, *output); err != nil {
		log.Fatal(fmt.Errorf("unable to build\n%w", err))
	}
}
