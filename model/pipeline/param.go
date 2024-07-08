package pipeline

type Param struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Default     string `yaml:"default,omitempty"`
}
