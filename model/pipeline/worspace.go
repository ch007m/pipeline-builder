package pipeline

type Workspace struct {
	Name     string `yaml:"name"`
	Optional bool   `yaml:"optional,omitempty"`
}
