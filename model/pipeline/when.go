package pipeline

type Operator string

const (
	In    Operator = "in"
	NotIn Operator = "notin"
)

type When struct {
	Input    string   `yaml:"input"`
	Operator Operator `yaml:"operator"`
	Values   []string `yaml:"values"`
}
