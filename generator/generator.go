package generator

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/logging"
	generator "github.com/ch007m/pipeline-builder/templates/lifecycle"

	"gopkg.in/yaml.v3"
)

func Contribute(path string, output string) error {
	configurator, err := NewConfigurator(path)
	if err != nil {
		return fmt.Errorf("Unable to read/parse the config yaml file %s, %w", path, err)
	}

	logging.Logger.Debug("Configurator path: %s", configurator)

	pipeline := generator.CreatePipeline()
	data, err := yaml.Marshal(&pipeline)
	if err != nil {
		return fmt.Errorf("Yaml marshalling error: %v\n", err)
	}

	return WriteFlow(data, &pipeline, output)
}
