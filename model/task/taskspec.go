package task

type TaskSpec struct {
	Description  string       `yaml:"description,omitempty"`
	Params       []Param      `yaml:"params,omitempty"`
	Results      []Result     `yaml:"results,omitempty"`
	Steps        []Step       `yaml:"steps,omitempty"`
	StepTemplate StepTemplate `yaml:"stepTemplate,omitempty"`
	Volumes      []Volume     `yaml:"volumes,omitempty"`
	Workspaces   []Workspace  `yaml:"workspaces,omitempty"`
}
