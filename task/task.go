package task

// Task represents a Tekton Task
type Task struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       TaskSpec `yaml:"spec"`
}
