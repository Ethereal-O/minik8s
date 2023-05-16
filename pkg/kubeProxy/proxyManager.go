package kubeProxy

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"sync"
)

var kubeProxyManager *KubeProxyManager
var Exited = make(chan bool)
var ToExit = make(chan bool)

func createKubeProxyManager() *KubeProxyManager {
	kubeProxyManager := &KubeProxyManager{}
	kubeProxyManager.RuntimeServiceMap = make(map[string]*object.RuntimeService)
	kubeProxyManager.RuntimeGatewayMap = make(map[string]*object.RuntimeGateway)
	var lock sync.Mutex
	kubeProxyManager.Lock = lock
	return kubeProxyManager
}

func Start_proxy() {
	fmt.Println("kube-proxy start")
	kubeProxyManager = createKubeProxyManager()
	kubeProxyManager.initKubeProxyManager()
	kubeProxyManager.startKubeProxyManager()
}

func (kubeProxyManager *KubeProxyManager) startKubeProxyManager() {
	runtimeServiceChan, runtimeServiceStop := messging.Watch("/"+config.RUNTIMESERVICE_TYPE, true)
	runtimeGatewayChan, runtimeGatewayStop := messging.Watch("/"+config.RUNTIMEGATEWAY_TYPE, true)
	go dealRuntimeService(runtimeServiceChan)
	go dealRuntimeGateway(runtimeGatewayChan)

	// Wait until Ctrl-C
	<-ToExit
	runtimeServiceStop()
	runtimeGatewayStop()
	Exited <- true
}

func dealRuntimeService(runtimeServiceChan chan string) {
	for {
		select {
		case mes := <-runtimeServiceChan:
			var tarRuntimeService object.RuntimeService
			err := json.Unmarshal([]byte(mes), &tarRuntimeService)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarRuntimeService.Service.Runtime.Status == config.EXIT_STATUS {
				dealExitRuntimeService(&tarRuntimeService)
			} else if tarRuntimeService.Service.Runtime.Status == config.RUNNING_STATUS {
				dealRunningRuntimeService(&tarRuntimeService)
			} else {
				fmt.Println("runtime service status error!")
			}
		}
	}
}

func dealRuntimeGateway(runtimeGatewayChan chan string) {
	for {
		select {
		case mes := <-runtimeGatewayChan:
			var tarRuntimeGateway object.RuntimeGateway
			err := json.Unmarshal([]byte(mes), &tarRuntimeGateway)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarRuntimeGateway.Gateway.Runtime.Status == config.EXIT_STATUS {
				dealExitRuntimeGateway(&tarRuntimeGateway)
			} else if tarRuntimeGateway.Gateway.Runtime.Status == config.RUNNING_STATUS {
				dealRunningRuntimeGateway(&tarRuntimeGateway)
			} else {
				fmt.Println("runtime gateway status error!")
			}
		}
	}
}
