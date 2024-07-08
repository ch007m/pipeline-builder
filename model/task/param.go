package task

// Param represents a parameter for a Task
type Param struct {
	Default     string `yaml:"default,omitempty"`
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
}
