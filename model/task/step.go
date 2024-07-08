package task

// Step represents a step in a Task
type Step struct {
	Name   string   `yaml:"name"`
	Image  string   `yaml:"image"`
	Env    []EnvVar `yaml:"env,omitempty"`
	Script string   `yaml:"script,omitempty"`
}
