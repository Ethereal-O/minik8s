package services

import (
	"minik8s/pkg/object"
	"sync"
	"time"
)

const (
	MAX_PODS = 5
)

// Service

type ServiceManager struct {
	StopChannel chan struct{}
	ServiceMap  map[string]ServiceStatus
	DnsMap      map[string]DnsStatus
	Lock        sync.Mutex
}

// Runtime service to check diff
type ServiceStatus struct {
	Service object.Service
	Pods    []object.Pod
	Timer   time.Ticker
	Status  string
	Lock    sync.Mutex
}

// Dns

type DnsStatus struct {
	Paths     []object.Path
	ClusterIp string
	Status    string
}
