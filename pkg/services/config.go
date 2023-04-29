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
	GATEWAY_STATUS_INIT                    = "init"
	GATEWAY_STATUS_DEPLOYING               = "deploying"
	GATEWAY_STATUS_RUNNING                 = "running"
	GATEWAY_PREFIX                         = "Gtw-"
	GATEWAY_REPLICASET_PREFIX              = "GtwRS-"
	GATEWAY_POD_PREFIX                     = "GtwPod-"
	GATEWAY_CONTAINER_PREFIX               = "GtwCont-"
	GATEWAY_SERVICE_PREFIX                 = "GtwSvc-"
	DNS_SERVICE_NAME                       = "DNS-Svc"
	DNS_GATEWAY_SELECTOR                   = "bound"
	NGINX_PATH_PREFIX                      = "/home/os/minik8s/Gateway"
	NGINX_CONFIG_FILE                      = "nginx.conf"
	HOST_PATH                              = "/home/os/minik8s/DNS/hosts.conf"
	NGINX_TEMPLATE_FILEPATH                = "template/config/NGINX_TEMPLATE"
	DNS_REPLICATESET_TEMPLATE_FILEPATH     = "template/yaml/dns-replicaset.yaml"
	DNS_SERVICE_TEMPLATE_FILEPATH          = "template/yaml/dns-service.yaml"
	GATEWAY_REPLICATESET_TEMPLATE_FILEPATH = "template/yaml/gateway-replicaset.yaml"
	GATEWAY_SERVICE_TEMPLATE_FILEPATH      = "template/yaml/gateway-service.yaml"
)

// Service

type ServiceManager struct {
	ServiceMap map[string]*object.RuntimeService
	Lock       sync.Mutex
}

// Gateway

type DnsManager struct {
	Timer              time.Ticker
	DnsTemplates       DnsTemplate
	GatewayMap         map[string]*object.RuntimeGateway
	ToBeDoneGatewayMap map[string]object.RuntimeGateway
	Lock               sync.Mutex
}

type DnsTemplate struct {
	DnsReplicaSetTemplate     object.ReplicaSet
	DnsServiceTemplate        object.Service
	GateWayReplicaSetTemplate object.ReplicaSet
	GateWayServiceTemplate    object.Service
}
