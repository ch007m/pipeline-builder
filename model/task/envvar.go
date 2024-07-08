package task

// EnvVar represents an environment variable for a Step
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
