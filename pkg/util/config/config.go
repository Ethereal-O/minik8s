package config

import (
	"io/ioutil"
)

var (
	// MASTER_IP is the IP address of master node
	MASTER_IP string
	LOCALHOST string

	// ETCD_ENDPOINT / NSQ_PRODUCER are used by only master node
	ETCD_ENDPOINT string
	NSQ_PRODUCER  string

	// APISERVER_URL / NSQ_CONSUMER are used by every worker node
	APISERVER_URL string
	NSQ_CONSUMER  string

	// DNS_SERVER makes all server have same DNS server to support DNS query
	DNS_SERVER string
)

const (
	POD_TYPE            = "Pod"
	REPLICASET_TYPE     = "Replicaset"
	SERVICE_TYPE        = "Service"
	RUNTIMESERVICE_TYPE = "RuntimeService"
	NODE_TYPE           = "Node"
	DNS_TYPE            = "DNS"
	GATEWAY_TYPE        = "Gateway"
	RUNTIMEGATEWAY_TYPE = "RuntimeGateway"
)

var TP = []string{POD_TYPE, REPLICASET_TYPE, SERVICE_TYPE, RUNTIMESERVICE_TYPE, NODE_TYPE, DNS_TYPE, GATEWAY_TYPE, RUNTIMEGATEWAY_TYPE}

const EMPTY_FLAG = "none"

const (
	// Only for pod
	CREATED_STATUS = "CREATED"
	BOUND_STATUS   = "BOUND"

	// For all API objects
	RUNNING_STATUS = "RUNNING"
	EXIT_STATUS    = "EXIT"
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

	APISERVER_URL = "http://" + MASTER_IP + ":8080"
	NSQ_CONSUMER = MASTER_IP + ":4161"

	DNS_SERVER = "11.11.11.11"
}
