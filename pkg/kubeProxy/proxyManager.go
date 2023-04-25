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
	kubeProxyManager.RootMap = make(map[string]map[string]*SingleService)
	kubeProxyManager.RuntimeServiceMap = make(map[string]object.RuntimeService)
	kubeProxyManager.GatewayMap = make(map[string]object.Gateway)
	var lock sync.Mutex
	kubeProxyManager.Lock = lock
	return kubeProxyManager
}

func Start_proxy() {
	fmt.Println("kube-proxy start")
	kubeProxyManager = createKubeProxyManager()
	kubeProxyManager.initRootChain()
	kubeProxyManager.initKubeProxyManager()
}

func (kubeProxyManager *KubeProxyManager) initKubeProxyManager() {
	runtimeServiceChan, runtimeServiceStop := messging.Watch("/"+config.RUNTIMESERVICE_TYPE, true)
	go dealRuntimeService(runtimeServiceChan)

	// Wait until Ctrl-C
	<-ToExit
	runtimeServiceStop()
	kubeProxyManager.deleteRootChain()
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
