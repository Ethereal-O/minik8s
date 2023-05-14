package config

import (
	"io/ioutil"
)

var (
	// MASTER_IP is the IP address of master node
	MASTER_IP string
	LOCALHOST string

	// ETCD_ENDPOINT / NSQ_PRODUCER / PROMETHEUS_URL are used by only master node
	ETCD_ENDPOINT  string
	NSQ_PRODUCER   string
	PROMETHEUS_URL string

	// APISERVER_URL / NSQ_CONSUMER are used by every worker node
	APISERVER_URL string
	NSQ_CONSUMER  string

	// DNS_SERVER makes all server have same DNS server to support DNS query
	DNS_SERVER string
)

const (
	POD_TYPE            = "Pod"
	REPLICASET_TYPE     = "ReplicaSet"
	AUTOSCALER_TYPE     = "AutoScaler"
	SERVICE_TYPE        = "Service"
	RUNTIMESERVICE_TYPE = "RuntimeService"
	NODE_TYPE           = "Node"
	DNS_TYPE            = "DNS"
	GATEWAY_TYPE        = "Gateway"
	RUNTIMEGATEWAY_TYPE = "RuntimeGateway"
	GPUJOB_TYPE         = "GpuJob"
	GPUFILE_TYPE        = "GpuFile"
)

var TP = []string{POD_TYPE, REPLICASET_TYPE, AUTOSCALER_TYPE, SERVICE_TYPE, RUNTIMESERVICE_TYPE, NODE_TYPE, DNS_TYPE, GATEWAY_TYPE, RUNTIMEGATEWAY_TYPE, GPUJOB_TYPE, GPUFILE_TYPE}

const EMPTY_FLAG = "none"

const (
	// Only for pod
	CREATED_STATUS = "CREATED"
	BOUND_STATUS   = "BOUND"

	// For all API objects
	RUNNING_STATUS = "RUNNING"
	EXIT_STATUS    = "EXIT"
)

const (
	SERVICE_TYPE_NODEPORT  = "NodePort"
	SERVICE_TYPE_CLUSTERIP = "ClusterIP"
	DNS_SERVICE_NAME       = "DNS-Svc"
)

// Some const for GPU
const (
	NODE_DIR_PATH      = "/home/shareDir"
	CONTAINER_DIR_PATH = "/home/shareDir"
	GPU_JOB_NAME       = "gpujob"
	GPU_JOB_IMAGE      = "henry35/zsh-gpu-server:4.0"
	GPU_JOB_COMMAND    = "/home/server.sh"
)

func init() {
	data, err := ioutil.ReadFile("master_ip.txt")
	if err != nil {
		panic(err)
	}

	MASTER_IP = string(data)
	LOCALHOST = "localhost"

	ETCD_ENDPOINT = LOCALHOST + ":2379"
	NSQ_PRODUCER = LOCALHOST + ":4150"
	PROMETHEUS_URL = "http://" + LOCALHOST + ":9090"

	APISERVER_URL = "http://" + MASTER_IP + ":8080"
	NSQ_CONSUMER = MASTER_IP + ":4161"

	DNS_SERVER = "10.10.10.10"
}
