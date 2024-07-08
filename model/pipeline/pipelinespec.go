package pipeline

type Spec struct {
	Workspaces []Workspace `yaml:"workspaces"`
	Params     []Param     `yaml:"params"`
	Results    []Result    `yaml:"results"`
	Finally    []Finally   `yaml:"finally"`
	Tasks      []Task      `yaml:"tasks"`
}
