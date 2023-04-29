package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"sync"
)

var serviceManager *ServiceManager
var Exited = make(chan bool)
var ToExit = make(chan bool)

func createServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	serviceManager.ServiceMap = make(map[string]*object.RuntimeService)
	var lock sync.Mutex
	serviceManager.Lock = lock
	return serviceManager
}

func StartServiceManager() {
	serviceManager = createServiceManager()
	serviceChan, serviceStop := messging.Watch("/"+config.SERVICE_TYPE, true)
	go dealService(serviceChan)

	// Wait until Ctrl-C
	<-ToExit
	serviceStop()
	Exited <- true
}

func dealService(serviceChan chan string) {
	for {
		select {
		case mes := <-serviceChan:
			if mes == "hello" {
				continue
			}
			var tarService object.Service
			err := json.Unmarshal([]byte(mes), &tarService)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarService.Runtime.Status == config.EXIT_STATUS {
				dealExitService(&tarService)
			} else if tarService.Runtime.Status == config.RUNNING_STATUS {
				dealRunningService(&tarService)
			} else {
				fmt.Println("service status error!")
			}
		}
	}
}
