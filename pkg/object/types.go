package object

// --------------------------- Basic Types ---------------------------

type Labels map[string]string

// Metadata take values from .yaml files
type Metadata struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
	Labels    Labels `yaml:"labels" json:"labels"`
}

// Runtime generate values from runtime (not in .yaml files)
type Runtime struct {
	Uuid   string `yaml:"uuid" json:"uuid"`
	Status string `yaml:"status" json:"status"`
	// When a pod belongs to a replica set, Belong refers to the Name of the replica set
	Belong string `yaml:"belong" json:"belong"`
	// When a pod is bound to a node, Bind refers to the Name of the node
	Bind string `yaml:"bind" json:"bind"`
	// When a pod is created, it should have a cluster IP so that the containers of the pod share network namespace
	ClusterIp string `yaml:"clusterIp" json:"clusterIp"`
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
	RestartPolicy string      `yaml:"restartPolicy" json:"restartPolicy"`
	Volumes       []Volume    `yaml:"volumes" json:"volumes"`
	Containers    []Container `yaml:"containers" json:"containers"`
	NodeSelector  Labels      `yaml:"nodeSelector" json:"nodeSelector"`
}

type Volume struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
	Path string `yaml:"path" json:"path"`
}

type Container struct {
	Name            string        `yaml:"name" json:"name"`
	Image           string        `yaml:"image" json:"image"`
	ImagePullPolicy string        `yaml:"imagePullPolicy" json:"imagePullPolicy"`
	Ports           []Port        `yaml:"ports" json:"ports"`
	VolumeMounts    []VolumeMount `yaml:"volumeMounts" json:"volumeMounts"`
	Limits          Limits        `yaml:"limits" json:"limits"`
	Args            []string      `yaml:"args" json:"args"`
	Command         []string      `yaml:"cmd" json:"cmd"`
	Env             []EnvVar      `yaml:"env"`
}

type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type Port struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	Protocol      string `yaml:"protocol" json:"protocol"`
	HostPort      int    `yaml:"hostPort" json:"hostPort"`
	HostIP        string `yaml:"hostIP" json:"hostIP"`
}

type VolumeMount struct {
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
	Ports    []Ports  `yaml:"ports" json:"ports"`
	Selector Selector `yaml:"selector" json:"selector"`
	Type     string   `yaml:"type" json:"type"`
	Template Template `yaml:"template" json:"template"`
	Strategy Strategy `yaml:"strategy" json:"strategy"`
	Replicas int      `yaml:"replicas" json:"replicas"`
}

type Selector struct {
	Name string `yaml:"name" json:"name"`
}

type Strategy struct {
	Type string `yaml:"type" json:"type"`
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
