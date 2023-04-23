package object

// --------------------------- Basic Types ---------------------------

// Metadata take values from .yaml files
type Metadata struct {
	Name   string            `yaml:"name" json:"name"`
	Labels map[string]string `yaml:"labels" json:"labels"`
}

// Runtime generate values from runtime (not in .yaml files)
type Runtime struct {
	Uuid   string `yaml:"uuid" json:"uuid"`
	Status string `yaml:"status" json:"status"`
	// When a pod belongs to a replica set, Belong refers to the Name of the replica set
	Belong string `yaml:"belong" json:"belong"`
	// When a pod is bound to a node, Bind refers to the Name of the node
	Bind string `yaml:"bind" json:"bind"`
}

// --------------------------- Node ---------------------------

type Node struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     NodeSpec `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
}

type NodeSpec struct {
	Ip string `yaml:"ip" json:"ip"`
}

// --------------------------- Replica Set ---------------------------

type ReplicaSet struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     RsSpec   `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
}

type RsSpec struct {
	Replicas int      `yaml:"replicas" json:"replicas"`
	Template Template `yaml:"template" json:"template"`
}

type Template struct {
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
}

// --------------------------- Pod ---------------------------

type Pod struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     PodSpec  `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
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
	Command      []string       `yaml:"cmd" json:"cmd"`
}

type Ports struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	Protocol      string `yaml:"protocol" json:"protocol"`
}

type VolumeMounts struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mountPath"`
}

type Limits struct {
	Cpu    string `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
}

// --------------------------- Service ---------------------------

type Service struct {
	Kind     string      `yaml:"kind" json:"kind"`
	Metadata Metadata    `yaml:"metadata" json:"metadata"`
	Spec     ServiceSpec `yaml:"spec" json:"spec"`
	Runtime  Runtime     `yaml:"runtime" json:"runtime"`
}

type ServiceSpec struct {
	Type      string            `json:"type" yaml:"type"`
	ClusterIp string            `json:"clusterIp" yaml:"clusterIp"`
	Ports     []ServicePort     `json:"ports" yaml:"ports"`
	Selector  map[string]string `json:"selector" yaml:"selector"`
}

type ServicePort struct {
	Name       string `json:"name" yaml:"name"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Port       string `json:"port" yaml:"port"`
	TargetPort string `json:"targetPort" yaml:"targetPort"`
	NodePort   string `json:"nodePort" yaml:"nodePort"`
}

// --------------------------- Dns ---------------------------

type Dns struct {
	Kind     string   `yaml:"kind" json:"kind"`
	MetaData Metadata `json:"metadata" yaml:"metadata"`
	Spec     DnsSpec  `json:"spec" yaml:"spec"`
}

type DnsSpec struct {
	Host  string `yaml:"host" json:"host"`
	Paths []Path `yaml:"paths" json:"paths"`
}

type Path struct {
	Name    string  `yaml:"name" json:"name"`
	Service Service `yaml:"service" json:"service"`
	IP      string  `yaml:"ip" json:"ip"`
	Port    string  `yaml:"port" json:"port"`
}
