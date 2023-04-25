package kubeProxy

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var kubeProxyManager *KubeProxyManager
var kubeProxyManagerExited = make(chan bool)
var kubeProxyManagerToExit = make(chan os.Signal)

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
	signal.Notify(kubeProxyManagerToExit, syscall.SIGINT, syscall.SIGTERM)
	runtimeServiceChan, runtimeServiceStop := messging.Watch("/"+config.RUNTIMESERVICE_TYPE, true)
	go dealRuntimeService(runtimeServiceChan)

	// Wait until Ctrl-C
	<-kubeProxyManagerToExit
	runtimeServiceStop()
	kubeProxyManager.deleteRootChain()
	kubeProxyManagerExited <- true
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
