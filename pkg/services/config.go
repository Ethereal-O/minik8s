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
	GATEWAY_PREFIX                         = "$gateway-"
	GATEWAY_REPLICASET_PREFIX              = "$gateWayReplicaset-"
	GATEWAY_POD_PREFIX                     = "gateWayPod-"
	GATEWAY_CONTAINER_PREFIX               = "gateWayContainer-"
	GATEWAY_SERVICE_PREFIX                 = "gateWayService-"
	DNS_GATEWAY_SELECTOR                   = "$bound$"
	NGINX_PATH_PREFIX                      = "/root/nginx"
	NGINX_CONFIG_FILE                      = "nginx.conf"
	HOST_PATH                              = "/root/dns/hosts.conf"
	DNS_REPLICATESET_TEMPLATE_FILEPATH     = "templates/dns-replicaset.yaml"
	DNS_SERVICE_TEMPLATE_FILEPATH          = "templates/dns-service.yaml"
	GATEWAY_REPLICATESET_TEMPLATE_FILEPATH = "templates/gateway-replicaset.yaml"
	GATEWAY_SERVICE_TEMPLATE_FILEPATH      = "templates/gateway-service.yaml"
)

// Service

type ServiceManager struct {
	ServiceMap map[string]object.ServiceStatus
	Lock       sync.Mutex
}

// Gateway

type DnsManager struct {
	Timer        time.Ticker
	DnsTemplates DnsTemplate
	GatewayMap   map[string]object.GatewayStatus
	Lock         sync.Mutex
}

type DnsTemplate struct {
	DnsReplicaSetTemplate     object.ReplicaSet
	DnsServiceTemplate        object.Service
	GateWayReplicaSetTemplate object.ReplicaSet
	GateWayServiceTemplate    object.Service
}
