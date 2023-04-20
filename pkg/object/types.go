package object

type ReplicaSet struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     RsSpec   `yaml:"spec" json:"spec"`
}

type RsSpec struct {
	Replicas int      `yaml:"replicas" json:"replicas"`
	Template Template `yaml:"template" json:"template"`
}

type Template struct {
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
}

// ------------------------------------------------

type Pod struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
	Belong   string   `yaml:"belong" json:"belong"`
}

type Metadata struct {
	Name   string            `yaml:"name" json:"name"`
	Labels map[string]string `yaml:"labels" json:"labels"`
	Uuid   string            `yaml:"uuid" json:"uuid"`
	Status string            `yaml:"status" json:"status"`
}

type PodSpec struct {
	Volumes    []Volumes    `yaml:"volumes" json:"volumes"`
	Containers []Containers `yaml:"containers" json:"containers"`
}

type Volumes struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
	Path string `yaml:"path" json:"path"`
}

type Containers struct {
	Name         string         `yaml:"name" json:"name"`
	Image        string         `yaml:"image" json:"image"`
	Ports        []Ports        `yaml:"ports" json:"ports"`
	VolumeMounts []VolumeMounts `yaml:"volumeMounts" json:"volumeMounts"`
	Limits       Limits         `yaml:"limits" json:"limits"`
	Args         []string       `yaml:"args" json:"args"`
}

type Ports struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	Protocol      string `yaml:"protocol" json:"protocol"`
}

type VolumeMounts struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mountPath"'`
}

type Limits struct {
	Cpu    string `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
}
