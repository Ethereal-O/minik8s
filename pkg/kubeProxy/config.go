package kubeProxy

import (
	"minik8s/pkg/object"
	"sync"
)

const (
	SINGLE_SERVICE   = "SINGLE_SERVICE"
	SINGLE_POD       = "SINGLE_POD"
	SINGLE_NET       = "SINGLE_NET"
	PARENT_TABLE     = "nat"
	PARENT_CHAIN     = "KUBE_PROXY_PARENT_CHAIN"
	OUTPUT_CHAIN     = "OUTPUT"
	PREROUTING_CHAIN = "PREROUTING"
)

type KubeProxyManager struct {
	// RootMap is a map of map, the first key is the service name, the second key is the pod port, and the value is a single service
	RootMap           map[string]map[string]*SingleService
	RuntimeServiceMap map[string]object.RuntimeService
	GatewayMap        map[string]object.Gateway
	Lock              sync.Mutex
}

type PodInfo struct {
	// PodName is the name of the pod
	PodName string
	// PodIP is the ip of the pod
	PodIP string
	// PodPort is the port of the pod
	PodPort string
}

type SingleService struct {
	Table        string
	Parent       string
	Name         string
	ClusterIp    string
	ClusterPort  string
	Protocol     string
	SinglePodMap map[string]*SinglePod
	RuleCommand  []string
}

type SinglePod struct {
	Table       string
	Parent      string
	Name        string
	PodName     string
	Protocol    string
	Net         *SingleNet
	RuleCommand []string
	Id          int
}

type SingleNet struct {
	Table       string
	Parent      string
	Name        string
	PodIp       string
	Port        string
	Protocol    string
	RuleCommand []string
}
