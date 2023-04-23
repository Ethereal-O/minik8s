package kubeProxy

import (
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/util/config"
	"sync"
)

func StartProxyManager() {
	fmt.Println("this is proxy")
}

var serviceManager *ServiceManager
var serviceManagerExited = make(chan bool)
var serviceManagerToExit = make(chan bool)

func createServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	serviceManager.ServiceMap = make(map[string]ServiceStatus)
	var lock sync.Mutex
	serviceManager.Lock = lock
	return serviceManager
}

func StartServiceManager() {
	serviceManager = createServiceManager()
	serviceChan, serviceStop := messging.Watch("/"+config.SERVICE_TYPE, true)
	go dealService(serviceChan)

	// Wait until Ctrl-C
	<-serviceManagerToExit
	serviceStop()
	serviceManagerExited <- true
}
