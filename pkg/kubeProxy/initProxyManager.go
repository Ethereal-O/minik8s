package kubeProxy

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/tools"
	"time"
)

func (kubeProxyManager *KubeProxyManager) initKubeProxyManager() {
	updateDnsConfig()
	kubeProxyManager.Timer = *time.NewTicker(CHECK_NODEPORT_SERVICE_TIME_INTERVAL)
	go kubeProxyManager.checkNodePortServiceLoop()
}

func (kubeProxyManager *KubeProxyManager) checkNodePortServiceLoop() {
	defer kubeProxyManager.Timer.Stop()
	for {
		select {
		case <-kubeProxyManager.Timer.C:
			kubeProxyManager.checkNodePortService()
		}
	}
}

func (kubeProxyManager *KubeProxyManager) checkNodePortService() {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	kubeProxyManager.updateRuntimeService()
}

func (kubeProxyManager *KubeProxyManager) updateRuntimeService() {
	// get all runtimeService
	allRuntimeServices := client.GetAllRuntimeServices()

	// first check if service is running
	runningRuntimeServices, _ := tools.Filter(allRuntimeServices, func(runtimeService object.RuntimeService) bool {
		if runtimeService.Status == services.SERVICE_STATUS_RUNNING && runtimeService.Service.Spec.Type == config.SERVICE_TYPE_NODEPORT {
			return true
		} else {
			return false
		}
	})

	if len(runningRuntimeServices) == 0 {
		return
	}

	// apply filter to get new runtimeService
	newRuntimeServices, _ := tools.Filter(runningRuntimeServices, func(runtimeService object.RuntimeService) bool {
		if _, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]; ok {
			return false
		} else {
			return true
		}
	})

	if len(newRuntimeServices) == 0 {
		return
	}

	fmt.Printf("updating node port service config num %d\n", len(newRuntimeServices))

	// add new runtimeService to map
	for _, runtimeService := range newRuntimeServices {
		runtimeServiceRef := runtimeService
		kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = &runtimeServiceRef
	}

	// update node port service config
	for _, runtimeService := range kubeProxyManager.RuntimeServiceMap {
		applyNodePortService(runtimeService)
	}
}
