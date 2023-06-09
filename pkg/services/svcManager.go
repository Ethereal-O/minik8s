package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"sync"
	"time"
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

	time.Sleep(5 * time.Second)
	enableForwarding()

	// Wait until Ctrl-C
	<-ToExit
	serviceStop()
	Exited <- true
}

func dealService(serviceChan chan string) {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_NGINX {
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

	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES || config.SERVICE_POLICY == config.SERVICE_POLICY_MICROSERVICE {
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
					dealExitService_old(&tarService)
				} else if tarService.Runtime.Status == config.RUNNING_STATUS {
					dealRunningService_old(&tarService)
				} else {
					fmt.Println("service status error!")
				}
			}
		}
	}

}

func enableForwarding() {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_NGINX {
		value, _, _ := exeFile.ReadYaml(FORWARD_DAEMONSET_TEMPLATE_FILEPATH)
		var DaemonSetObject object.DaemonSet
		err := json.Unmarshal([]byte(value), &DaemonSetObject)
		if err != nil {
			fmt.Println("Enable forwarding fail" + err.Error())
			return
		}
		client.AddDaemonSet(DaemonSetObject)
	}
}
