package config

const (
	//BASE_URL              = "192.168.29.132"
	BASE_URL              = "192.168.31.68"
	APISERVER_URL  string = "http://" + BASE_URL + ":8080"
	ETCD_Endpoints string = BASE_URL + ":2379"
	NSQ_PEODUCER   string = BASE_URL + ":4150"
	NSQ_CONSUMER   string = BASE_URL + ":4161"
	// make all server have same DNS server to support DNS query
	DNS_SERVER string = "11.11.11.11"
)

const (
	POD_TYPE           = "Pod"
	REPLICASET_TYPE    = "Replicaset"
	SERVICE_TYPE       = "Service"
	SERVICESTATUS_TYPE = "ServiceStatus"
	NODE_TYPE          = "Node"
	DNS_TYPE           = "DNS"
	GATEWAY_TYPE       = "Gateway"
	GATEWAYSTATUS_TYPE = "GatewayStatus"
)

var TP = []string{POD_TYPE, REPLICASET_TYPE, SERVICE_TYPE, SERVICESTATUS_TYPE, NODE_TYPE, DNS_TYPE, GATEWAY_TYPE, GATEWAYSTATUS_TYPE}

const EMPTY_FLAG = "none"

const (
	// Only for pod
	CREATED_STATUS = "CREATED"
	BOUND_STATUS   = "BOUND"

	// For all API objects
	RUNNING_STATUS = "RUNNING"
	EXIT_STATUS    = "EXIT"
)
