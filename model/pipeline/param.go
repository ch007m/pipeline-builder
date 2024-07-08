package pipeline

type Param struct {
	Default     interface{} `yaml:"default,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type,omitempty"`
	Value       interface{} `yaml:"value,omitempty"`
}
