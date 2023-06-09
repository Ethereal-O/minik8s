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
	kubeProxyManager.VirtualServiceMap = make(map[string]*object.VirtualService)
	var lock sync.Mutex
	kubeProxyManager.Lock = lock
	return kubeProxyManager
}

func Start_proxy() {
	fmt.Println("kube-proxy start")
	kubeProxyManager = createKubeProxyManager()
	initialize()
	kubeProxyManager.startKubeProxyManager()
}

func (kubeProxyManager *KubeProxyManager) startKubeProxyManager() {
	runtimeServiceChan, runtimeServiceStop := messging.Watch("/"+config.RUNTIMESERVICE_TYPE, true)
	runtimeGatewayChan, runtimeGatewayStop := messging.Watch("/"+config.RUNTIMEGATEWAY_TYPE, true)
	virtualServiceChan, virtualServiceStop := messging.Watch("/"+config.VIRTUALSERVICE_TYPE, true)
	go dealRuntimeService(runtimeServiceChan)
	go dealRuntimeGateway(runtimeGatewayChan)
	go dealVirtualService(virtualServiceChan)

	// Wait until Ctrl-C
	<-ToExit
	finalize()
	runtimeServiceStop()
	runtimeGatewayStop()
	virtualServiceStop()
	Exited <- true
}

func dealRuntimeService(runtimeServiceChan chan string) {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_NGINX {
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
	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES {
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
					dealExitRuntimeService_old(&tarRuntimeService)
				} else if tarRuntimeService.Service.Runtime.Status == config.RUNNING_STATUS {
					dealRunningRuntimeService_old(&tarRuntimeService)
				} else {
					fmt.Println("runtime service status error!")
				}
			}
		}
	}
	if config.SERVICE_POLICY == config.SERVICE_POLICY_MICROSERVICE {
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
					dealExitRuntimeService_micro(&tarRuntimeService)
				} else if tarRuntimeService.Service.Runtime.Status == config.RUNNING_STATUS {
					dealRunningRuntimeService_micro(&tarRuntimeService)
				} else {
					fmt.Println("runtime service status error!")
				}
			}
		}
	}
}

func initialize() {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES {
		kubeProxyManager.RootMap = make(map[string]map[string]*SingleService)
		kubeProxyManager.initRootChain()
	}
	if config.SERVICE_POLICY == config.SERVICE_POLICY_MICROSERVICE {
		kubeProxyManager.PodMatchMap = make(map[string]map[string]*PodMatch)
		kubeProxyManager.initSidecar()
	}
	kubeProxyManager.initKubeProxyManager()
}

func finalize() {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES {
		kubeProxyManager.deleteRootChain()
	}
	if config.SERVICE_POLICY == config.SERVICE_POLICY_MICROSERVICE {
		kubeProxyManager.deleteSidecar()
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

func dealVirtualService(virtualServiceChan chan string) {
	if config.SERVICE_POLICY != config.SERVICE_POLICY_MICROSERVICE {
		fmt.Println("You should use microservice policy to apply a virtual service!")
		return
	}
	for {
		select {
		case mes := <-virtualServiceChan:
			if mes == "hello" {
				continue
			}
			var tarVirtualService object.VirtualService
			err := json.Unmarshal([]byte(mes), &tarVirtualService)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarVirtualService.Runtime.Status == config.EXIT_STATUS {
				dealExitVirtualService(&tarVirtualService)
			} else if tarVirtualService.Runtime.Status == config.RUNNING_STATUS {
				dealRunningVirtualService(&tarVirtualService)
			} else {
				fmt.Println("virtual service status error!")
			}
		}
	}
}
