package task

// Param represents a parameter for a Task
type Param struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Default     string `yaml:"default,omitempty"`
}
