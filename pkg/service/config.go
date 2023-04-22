package service

import (
	"minik8s/pkg/object"
	"sync"
	"time"
)

// Service

type ServiceManager struct {
	stopChannel chan struct{}
	ServiceMap  map[string]*ServiceStatus
	DnsMap      map[string]*DnsStatus
	lock        sync.Mutex
}

// Runtime service to check diff
type ServiceStatus struct {
	Pods   []*object.Pod
	Timer  *time.Ticker
	Status string
	Error  error
}

// Dns

type DnsStatus struct {
	Paths     []*object.Path
	ClusterIp string
	Status    string
}
