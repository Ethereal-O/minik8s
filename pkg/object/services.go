package object

import (
	"sync"
	"time"
)

// Runtime service to check diff

type RuntimeService struct {
	Service Service
	Pods    []Pod
	Status  string
	Timer   time.Ticker `json:"-"`
	Lock    sync.Mutex  `json:"-"`
}

// Runtime gateway to deploy

type RuntimeGateway struct {
	Gateway   Gateway
	ClusterIp string
	Status    string
}
