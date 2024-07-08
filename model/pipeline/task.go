package pipeline

// Task represents a Tekton Task part of a Pipeline
type Task struct {
	Name       string             `yaml:"name"`
	RunAfter   []string           `yaml:"runAfter"`
	When       []When             `yaml:"when"`
	Params     []Param            `yaml:"params"`
	TaskRef    TaskRef            `yaml:"taskRef"`
	Workspaces []WorkspaceBinding `yaml:"workspaces"`
}
