package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
)

func dealRunningRuntimeService_micro(runtimeService *object.RuntimeService) {
	oldRuntimeService, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		fmt.Printf("creating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		createRuntimeService_micro(runtimeService)
	} else if tools.MD5(oldRuntimeService.Service) != tools.MD5(*runtimeService) {
		fmt.Printf("updating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		updateRuntimeService_micro(runtimeService)
	} else {
		fmt.Printf("duplicated runtimeService %s\n", runtimeService.Service.Metadata.Name)
	}
}

func dealExitRuntimeService_micro(runtimeService *object.RuntimeService) {
	fmt.Printf("deleting runtimeService %s\n", runtimeService.Service.Metadata.Name)
	deleteRuntimeService_micro(runtimeService)
}

func createRuntimeService_micro(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService
	// update podMatchMap
	createMultiService(runtimeService)
}

func deleteRuntimeService_micro(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	runtimeService, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		return
	}
	delete(kubeProxyManager.PodMatchMap, runtimeService.Service.Runtime.ClusterIp)
	delete(kubeProxyManager.RuntimeServiceMap, runtimeService.Service.Metadata.Name)
}

func updateRuntimeService_micro(runtimeService *object.RuntimeService) {
	deleteRuntimeService_micro(runtimeService)
	createRuntimeService_micro(runtimeService)
}
