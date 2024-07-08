package task

// Result represents a result produced by a Task
type Result struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}
