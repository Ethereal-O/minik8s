package kubeProxy

import (
	"minik8s/pkg/object"
	"sync"
)

const (
	SINGLE_SERVICE    = "Svc"
	SINGLE_POD        = "POD"
	SINGLE_NET        = "NET"
	ROOT_TABLE        = "nat"
	ROOT_CHAIN        = "KUBE_PROXY_PARENT_CHAIN"
	OUTPUT_CHAIN      = "OUTPUT"
	PREROUTING_CHAIN  = "PREROUTING"
	POSTROUTING_CHAIN = "POSTROUTING"
)

type KubeProxyManager struct {
	// RootMap is a map of map, the first key is the service name, the second key is the pod port, and the value is a single service
	RootMap           map[string]map[string]*SingleService
	RuntimeServiceMap map[string]*object.RuntimeService
	RuntimeGatewayMap map[string]*object.RuntimeGateway
	RootChain         RootChain
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

type RootChain struct {
	RuleCommand []string
}

type SingleService struct {
	Table                string
	Parent               string
	Name                 string
	ClusterIp            string
	ClusterPort          string
	IsNodePort           bool
	NodePort             string
	Protocol             string
	SinglePodMap         map[string]*SinglePod
	RuleCommandClusterIp []string
	RuleCommandNodePort  []string
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
