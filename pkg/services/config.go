package services

import (
	"minik8s/pkg/object"
	"sync"
	"time"
)

const (
	MAX_PODS                               = 5
	CHECK_PODS_TIME_INTERVAL               = 5 * time.Second
	CHECK_DNS_TIME_INTERVAL                = 5 * time.Second
	CREATE_RS_AND_SERVICE_TIME_INTERVAL    = 1 * time.Second
	UPDATE_SERVICE_TIME_INTERVAL           = 1 * time.Second
	SERVICE_STATUS_INIT                    = "INIT"
	SERVICE_STATUS_RUNNING                 = "RUNNING"
	SERVICE_STATUS_EXIT                    = "EXIT"
	GATEWAY_STATUS_INIT                    = "INIT"
	GATEWAY_STATUS_DEPLOYING               = "DEPLOYING"
	GATEWAY_STATUS_RUNNING                 = "RUNNING"
	GATEWAY_PREFIX                         = "Gtw-"
	GATEWAY_REPLICASET_PREFIX              = "GtwRS-"
	GATEWAY_POD_PREFIX                     = "GtwPod-"
	GATEWAY_CONTAINER_PREFIX               = "GtwCont-"
	GATEWAY_SERVICE_PREFIX                 = "GtwSvc-"
	DNS_SERVICE_NAME                       = "DNS-Svc"
	SERVICE_REPLICASET_PREFIX              = "SvcRS-"
	SERVICE_POD_PREFIX                     = "SvcPod-"
	SERVICE_CONTAINER_PREFIX               = "SvcCont-"
	FORWARD_DAEMONSET_PREFIX               = "Forward-DS"
	ALL_SELECTOR                           = "bound"
	GATEWAY_NGINX_PATH_PREFIX              = "/home/os/minik8s/Gateway"
	SERVICE_NGINX_PATH_PREFIX              = "/home/os/minik8s/Service"
	FORWARD_NGINX_PATH                     = "/home/os/minik8s/Forward/nginx.conf"
	NGINX_CONFIG_FILE                      = "nginx.conf"
	HOST_PATH                              = "/home/os/minik8s/DNS/hosts.conf"
	HOST_HOSTS_PATH                        = "/etc/hosts"
	HOST_HOSTS_BAK_PATH                    = "/etc/hosts.bak"
	NGINX_TEMPLATE_FILEPATH                = "template/config/NGINX_TEMPLATE"
	SERVICE_REPLICATESET_TEMPLATE_FILEPATH = "template/yaml/service-replicaset.yaml"
	DNS_REPLICATESET_TEMPLATE_FILEPATH     = "template/yaml/dns-replicaset.yaml"
	DNS_SERVICE_TEMPLATE_FILEPATH          = "template/yaml/dns-service.yaml"
	GATEWAY_REPLICATESET_TEMPLATE_FILEPATH = "template/yaml/gateway-replicaset.yaml"
	GATEWAY_SERVICE_TEMPLATE_FILEPATH      = "template/yaml/gateway-service.yaml"
	FORWARD_DAEMONSET_TEMPLATE_FILEPATH    = "template/yaml/forward-daemonset.yaml"
)

// Service

type ServiceManager struct {
	ServiceMap map[string]*object.RuntimeService
	Lock       sync.Mutex
}

// Gateway

type DnsManager struct {
	Timer              time.Ticker
	Templates          Template
	GatewayMap         map[string]*object.RuntimeGateway
	ToBeDoneGatewayMap map[string]object.RuntimeGateway
	Lock               sync.Mutex
}

type Template struct {
	DnsReplicaSetTemplate     object.ReplicaSet
	DnsServiceTemplate        object.Service
	GateWayReplicaSetTemplate object.ReplicaSet
	GateWayServiceTemplate    object.Service
	ServiceReplicaSetTemplate object.ReplicaSet
}
