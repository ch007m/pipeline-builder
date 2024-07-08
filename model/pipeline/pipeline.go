package pipeline

import "github.com/ch007m/pipeline-builder/model/common"

// Pipeline represents a Tekton Pipeline
type Pipeline struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   common.Metadata `yaml:"metadata"`
	Spec       Spec            `yaml:"spec"`
}
