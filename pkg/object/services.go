package object

import (
	"sync"
	"time"
)

// Runtime service to check diff

type ServiceStatus struct {
	Service Service
	Pods    []Pod
	Timer   time.Ticker
	Lock    sync.Mutex
}

// Runtime gateway to deploy

type GatewayStatus struct {
	Gateway   Gateway
	ClusterIp string
	Status    string
}
