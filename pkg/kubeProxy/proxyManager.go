package kubeProxy

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/iptables"
	"sync"
)

var kubeProxyManager *KubeProxyManager
var kubeProxyManagerExited = make(chan bool)
var kubeProxyManagerToExit = make(chan bool)

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
	initParentChain()
	kubeProxyManager = createKubeProxyManager()
	kubeProxyManager.initKubeProxyManager()
}

func (kubeProxyManager *KubeProxyManager) initKubeProxyManager() {
	runtimeServiceChan, runtimeServiceStop := messging.Watch("/"+config.RUNTIMESERVICE_TYPE, true)
	go dealRuntimeService(runtimeServiceChan)

	// Wait until Ctrl-C
	<-kubeProxyManagerToExit
	runtimeServiceStop()
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
				fmt.Println("runtime service status error")
			}
		}
	}
}

func initParentChain() {
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("[chain] Boot error")
		fmt.Println(err)
	}
	exist, err2 := ipt.ChainExists(PARENT_TABLE, PARENT_CHAIN)
	if err2 != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	if exist {
		return
	}
	err = ipt.NewChain(PARENT_TABLE, PARENT_CHAIN)
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
		return
	}
	err = ipt.Insert(PARENT_TABLE, OUTPUT_CHAIN, 1, "-j", PARENT_CHAIN, "-s", "0/0", "-d", "0/0", "-p", "all")
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
		return
	}
	err = ipt.Insert(PARENT_TABLE, PREROUTING_CHAIN, 1, "-j", PARENT_CHAIN, "-s", "0/0", "-d", "0/0", "-p", "all")
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
		return
	}
}
