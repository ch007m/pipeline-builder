package task

import "github.com/ch007m/pipeline-builder/model/common"

// Task represents a Tekton Task
type Task struct {
	APIVersion   string          `yaml:"apiVersion"`
	Kind         string          `yaml:"kind"`
	Metadata     common.Metadata `yaml:"metadata"`
	Spec         TaskSpec        `yaml:"spec"`
	StepTemplate StepTemplate    `yaml:"stepTemplate"`
}
