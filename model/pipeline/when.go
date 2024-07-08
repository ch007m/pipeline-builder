package pipeline

type When struct {
	Input    string   `yaml:"input"`
	Operator string   `yaml:"operator"`
	Values   []string `yaml:"values"`
}
