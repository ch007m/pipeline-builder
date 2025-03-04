package pipeline

type Finally struct {
	Name       string             `yaml:"name"`
	When       []When             `yaml:"when"`
	Params     []Param            `yaml:"params"`
	TaskRef    TaskRef            `yaml:"taskRef"`
	Workspaces []WorkspaceBinding `yaml:"workspaces"`
}
