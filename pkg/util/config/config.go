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

	// APISERVER_URL / NSQ_CONSUMER / FUNCTION_PROXY_URL are used by every worker node
	APISERVER_URL      string
	FUNCTION_PROXY_URL string
	NSQ_CONSUMER       string

	// DNS_SERVER makes all server have same DNS server to support DNS query
	DNS_SERVER string
)

const (
	POD_TYPE                 = "Pod"
	REPLICASET_TYPE          = "ReplicaSet"
	DAEMONSET_TYPE           = "DaemonSet"
	AUTOSCALER_TYPE          = "AutoScaler"
	SERVICE_TYPE             = "Service"
	RUNTIMESERVICE_TYPE      = "RuntimeService"
	VIRTUALSERVICE_TYPE      = "VirtualService"
	NODE_TYPE                = "Node"
	DNS_TYPE                 = "DNS"
	GATEWAY_TYPE             = "Gateway"
	RUNTIMEGATEWAY_TYPE      = "RuntimeGateway"
	GPUJOB_TYPE              = "GpuJob"
	SERVERLESSFUNCTIONS_TYPE = "ServerlessFunctions"
	FUNCTION_TYPE            = "Function"
	TRANSFILE_TYPE           = "TransFile"
)

var TP = []string{POD_TYPE, REPLICASET_TYPE, DAEMONSET_TYPE, AUTOSCALER_TYPE, SERVICE_TYPE, RUNTIMESERVICE_TYPE, VIRTUALSERVICE_TYPE, NODE_TYPE, DNS_TYPE, GATEWAY_TYPE, RUNTIMEGATEWAY_TYPE,
	GPUJOB_TYPE, SERVERLESSFUNCTIONS_TYPE, FUNCTION_TYPE, TRANSFILE_TYPE}

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
	SERVICE_TYPE_NODEPORT        = "NodePort"
	SERVICE_TYPE_CLUSTERIP       = "ClusterIP"
	SERVICE_POLICY               = SERVICE_POLICY_NGINX
	SERVICE_POLICY_NGINX         = "SERVICE_POLICY_NGINX"
	SERVICE_POLICY_IPTABLES      = "SERVICE_POLICY_IPTABLES"
	SERVICE_POLICY_MICROSERVICE  = "SERVICE_POLICY_MICROSERVICE"
	DNS_SERVICE_NAME             = "DNS-Svc"
	VIRTUAL_SERVICE_TYPE_EXACT   = "Exact"
	VIRTUAL_SERVICE_TYPE_PREFIX  = "Prefix"
	VIRTUAL_SERVICE_TYPE_REGULAR = "Regular"
)

// Some const for GPU
const (
	GPU_NODE_DIR_PATH      = "/home/shareDir"
	GPU_CONTAINER_DIR_PATH = "/home/shareDir"
	GPU_JOB_NAME           = "gpujob"
	GPU_JOB_IMAGE          = "henry35/zsh-gpu-server:4.0"
	GPU_JOB_COMMAND        = "/home/server.sh"
)

// Some const for Serverless
const (
	FUNC_NODE_DIR_PATH      = "/home/functions"
	FUNC_CONTAINER_DIR_PATH = "/home/functions"
	FUNC_NAME               = "func"
	FUNC_IMAGE              = "henry35/serverless:2.0"
	FUNC_COMMAND            = "/home/import.sh"
)

// Image source
const (
	PIP3_SOURCE_IMAGE_HOSTNAME = "pypi.tuna.tsinghua.edu.cn"
	PIP3_SOURCE_IMAGE_IP       = "101.6.15.130"
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
	FUNCTION_PROXY_URL = "http://" + MASTER_IP + ":8081"
	NSQ_CONSUMER = MASTER_IP + ":4161"

	if SERVICE_POLICY == SERVICE_POLICY_NGINX {
		DNS_SERVER = "10.10.10.10"
	}
	if SERVICE_POLICY == SERVICE_POLICY_IPTABLES {
		DNS_SERVER = "100.100.100.100"
	}
	if SERVICE_POLICY == SERVICE_POLICY_MICROSERVICE {
		DNS_SERVER = "114.114.114.114"
	}
}
