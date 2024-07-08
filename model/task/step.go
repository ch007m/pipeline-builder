package task

// Step represents a step in a Task
type Step struct {
	Name            string        `yaml:"name"`
	Image           string        `yaml:"image"`
	ImagePullPolicy string        `yaml:"imagePullPolicy"`
	Command         []string      `yaml:"command"`
	Args            []string      `yaml:"args"`
	Env             []EnvVar      `yaml:"env,omitempty"`
	Script          string        `yaml:"script,omitempty"`
	VolumeMounts    []VolumeMount `yaml:"volumeMount"`
}
