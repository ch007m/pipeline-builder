package task

type TaskRef struct {
	Resolver string  `yaml:"resolver"`
	Params   []Param `yaml:"params"`
}
