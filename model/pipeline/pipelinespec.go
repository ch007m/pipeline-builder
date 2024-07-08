package pipeline

import "github.com/ch007m/pipeline-builder/model/task"

type Spec struct {
	Workspaces []Workspace `yaml:"workspaces"`
	Params     []Param     `yaml:"params"`
	Results    []Result    `yaml:"results"`
	Finally    []Finally   `yaml:"finally"`
	Tasks      []task.Task `yaml:"tasks"`
}
