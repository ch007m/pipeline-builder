package task

type Workspace struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Optional    bool   `yaml:"optional,omitempty"`
}
