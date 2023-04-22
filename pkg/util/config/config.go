package config

const (
	BASE_URL              = "192.168.29.128"
	APISERVER_URL  string = "http://" + BASE_URL + ":8080"
	ETCD_Endpoints string = BASE_URL + ":2379"
	NSQ_PEODUCER   string = BASE_URL + ":4150"
	NSQ_CONSUMER   string = BASE_URL + ":4161"
)

const (
	POD_TYPE        = "Pod"
	REPLICASET_TYPE = "Replicaset"
	NODE_TYPE       = "Node"
)

var TP = []string{POD_TYPE, REPLICASET_TYPE, NODE_TYPE}

const EMPTY_FLAG = "none"

const (
	// Only for pod
	CREATED_STATUS = "CREATED"
	BOUND_STATUS   = "BOUND"

	// For all API objects
	RUNNING_STATUS = "RUNNING"
	EXIT_STATUS    = "EXIT"
)
