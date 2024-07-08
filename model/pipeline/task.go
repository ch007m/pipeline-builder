package pipeline

// Task represents a Tekton Task part of a Pipeline
type Task struct {
	Name       string             `yaml:"name"`
	RunAfter   []string           `yaml:"runAfter,omitempty"`
	When       []When             `yaml:"when,omitempty"`
	Params     []Param            `yaml:"params"`
	TaskRef    TaskRef            `yaml:"taskRef,omitempty"`
	Workspaces []WorkspaceBinding `yaml:"workspaces,omitempty"`
}
