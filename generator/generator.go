package generator

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/logging"
	"github.com/ch007m/pipeline-builder/templates/lifecycle"
	"github.com/ch007m/pipeline-builder/util"

	"gopkg.in/yaml.v3"
)

func Contribute(path string, output string) error {
	configurator, err := NewConfigurator(path)
	if err != nil {
		return fmt.Errorf("Unable to read/parse the config yaml file %s, %w", path, err)
	}

	logging.Logger.Debug("Configurator path: %s", configurator)

	pipeline := lifecycle.CreatePipeline()
	data, err := yaml.Marshal(&pipeline)
	if err != nil {
		return fmt.Errorf("Yaml marshalling error: %v\n", err)
	}

	return util.WriteFlow(data, &pipeline, output)
}
