package object

// --------------------------- Basic Types ---------------------------

type Labels map[string]string

// Metadata take values from .yaml files
type Metadata struct {
	Name   string `yaml:"name" json:"name"`
	Labels Labels `yaml:"labels" json:"labels"`
}

// Runtime generate values from runtime (not in .yaml files)
type Runtime struct {
	// --- Common ---
	Uuid   string `yaml:"uuid" json:"uuid"`
	Status string `yaml:"status" json:"status"`

	// --- Pod ---
	// When a pod belongs to a replica set, Belong refers to the Name of the replica set
	Belong string `yaml:"belong" json:"belong"`
	// When a pod is bound to a node, Bind refers to the Name of the node
	Bind string `yaml:"bind" json:"bind"`
	// When a pod is created, it should have a cluster IP so that the containers of the pod share network namespace
	ClusterIp string `yaml:"clusterIp" json:"clusterIp"`
	// PodIp is the IP address of the pod in docker (172.xx.xx.xx)
	PodIp string `yaml:"podIp" json:"podIp"`
	// When a pod is created, it should have a list of container ID
	Containers []string `yaml:"containers" json:"containers"`
	// Whether the pod should be restarted
	NeedRestart bool `yaml:"needRestart" json:"needRestart"`

	// --- ServerlessFunctions ---
	// When a function is available, FunctionIp is the target Ip
	FunctionIp string `yaml:"functionIp" json:"functionIp"`
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

// --------------------------- Daemon Set ---------------------------

type DaemonSet struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     DsSpec   `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
}

type DsSpec struct {
	Template Template `yaml:"template" json:"template"`
}

// --------------------------- Auto Scaler ---------------------------

// AutoScaler only support the ReplicaSet type currently
type AutoScaler struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     HpaSpec  `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
}

type HpaSpec struct {
	MinReplicas                       int       `yaml:"minReplicas" json:"minReplicas"`
	MaxReplicas                       int       `yaml:"maxReplicas" json:"maxReplicas"`
	Interval                          int       `yaml:"interval" json:"interval"`
	ScaleTargetRef                    TargetRef `yaml:"scaleTargetRef" json:"scaleTargetRef"`
	TargetCPUUtilizationPercentage    int       `yaml:"targetCPUUtilizationPercentage" json:"targetCPUUtilizationPercentage"`
	TargetCPUUtilizationStrategy      string    `yaml:"targetCPUUtilizationStrategy" json:"targetCPUUtilizationStrategy"`
	TargetMemoryUtilizationPercentage int       `yaml:"targetMemoryUtilizationPercentage" json:"targetMemoryUtilizationPercentage"`
	TargetMemoryUtilizationStrategy   string    `yaml:"targetMemoryUtilizationStrategy" json:"targetMemoryUtilizationStrategy"`
}

// TargetRef Kind only support ReplicaSet type currently
type TargetRef struct {
	Kind string `yaml:"kind" json:"kind"`
	Name string `yaml:"name" json:"name"`
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
	HostMode      string      `yaml:"hostMode" json:"hostMode"`
	Volumes       []Volume    `yaml:"volumes" json:"volumes"`
	Containers    []Container `yaml:"containers" json:"containers"`
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
	Type     string            `json:"type" yaml:"type"`
	Ports    []ServicePort     `json:"ports" yaml:"ports"`
	Selector map[string]string `json:"selector" yaml:"selector"`
}

type ServicePort struct {
	Name       string `json:"name" yaml:"name"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Port       string `json:"port" yaml:"port"`
	TargetPort string `json:"targetPort" yaml:"targetPort"`
	NodePort   string `json:"nodePort" yaml:"nodePort"`
}

// --------------------------- Gateway ---------------------------

type Gateway struct {
	Kind     string      `yaml:"kind" json:"kind"`
	Metadata Metadata    `json:"metadata" yaml:"metadata"`
	Spec     GatewaySpec `json:"spec" yaml:"spec"`
	Runtime  Runtime     `yaml:"runtime" json:"runtime"`
}

type GatewaySpec struct {
	Host  string `yaml:"host" json:"host"`
	Paths []Path `yaml:"paths" json:"paths"`
}

type Path struct {
	Name    string `yaml:"name" json:"name"`
	Service string `yaml:"service" json:"service"`
	IP      string `yaml:"ip" json:"ip"`
	Port    string `yaml:"port" json:"port"`
}

// --------------------------- GpuJob ---------------------------

type GpuJob struct {
	Kind     string     `yaml:"kind" json:"kind"`
	Metadata Metadata   `yaml:"metadata" json:"metadata"`
	Spec     GpuJobSpec `yaml:"spec" json:"spec"`
	Runtime  Runtime    `yaml:"runtime" json:"runtime"`
}

type GpuJobSpec struct {
	Path string `json:"path" yaml:"path"`
}

//----------------------------- ServerlessFunctions ---------------------------

type ServerlessFunctions struct {
	Kind     string   `yaml:"kind" json:"kind"`
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     FuncSpec `yaml:"spec" json:"spec"`
	Runtime  Runtime  `yaml:"runtime" json:"runtime"`
}

type FuncSpec struct {
	Path  string     `json:"path" yaml:"path"`
	Items []Function `json:"items" yaml:"items"`
}

type Function struct {
	FuncName string  `json:"funcName" yaml:"funcName"`
	Module   string  `json:"module" yaml:"module"`
	Runtime  Runtime `yaml:"runtime" json:"runtime"`
	FaasName string  `yaml:"faasName" json:"faasName"`
}

// --------------------------- File ---------------------------

type TransFile struct {
	Dirname string `json:"dirname" yaml:"dirname"`
	Data    string `json:"data" yaml:"data"`
	Tp      string `json:"tp" yaml:"tp"`
}

// --------------------------- DAG -----------------------------

type WorkFlow struct {
	WorkFlowName string    `json:"workFlowName" yaml:"workFlowName"`
	StartNode    string    `json:"startNode" yaml:"startNode"`
	Nodes        []DagNode `json:"nodes" yaml:"nodes"`
}

type DagNode struct {
	NodeName string   `json:"nodeName" yaml:"nodeName"`
	FuncName string   `json:"funcName" yaml:"funcName"`
	Choices  []Choice `json:"choices" yaml:"choices"`
}

type Choice struct {
	Condition Condition `json:"condition" yaml:"condition"`
	NextNode  string    `json:"nextNode" yaml:"nextNode"`
}

type Condition struct {
	TarVariable string `json:"tarVariable" yaml:"tarVariable"`
	TarValue    string `json:"tarValue" yaml:"tarValue"`
	Relation    string `json:"relation" yaml:"relation"`
}
