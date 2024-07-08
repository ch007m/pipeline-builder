package pipeline

import "github.com/ch007m/pipeline-builder/model/task"

type Finally struct {
	Name       string             `yaml:"name"`
	When       []When             `yaml:"when"`
	Params     []Param            `yaml:"params"`
	TaskRef    task.TaskRef       `yaml:"taskRef"`
	Workspaces []WorkspaceBinding `yaml:"workspaces"`
}
