package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
)

func dealRunningRuntimeService_old(runtimeService *object.RuntimeService) {
	oldRuntimeService, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		fmt.Printf("creating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		createRuntimeService_old(runtimeService)
	} else if tools.MD5(oldRuntimeService.Service) != tools.MD5(*runtimeService) {
		fmt.Printf("updating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		updateRuntimeService_old(runtimeService)
	} else {
		fmt.Printf("duplicated runtimeService %s\n", runtimeService.Service.Metadata.Name)
	}
}

func dealExitRuntimeService_old(runtimeService *object.RuntimeService) {
	fmt.Printf("deleting runtimeService %s\n", runtimeService.Service.Metadata.Name)
	deleteRuntimeService_old(runtimeService)
}

func createRuntimeService_old(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	// here create network
	fmt.Printf("creating network for runtimeService %s at %s!\n", runtimeService.Service.Metadata.Name, runtimeService.Service.Runtime.ClusterIp)
	createMultiService(runtimeService)
	fmt.Printf("creating network done for runtimeService %s at %s!\n", runtimeService.Service.Metadata.Name, runtimeService.Service.Runtime.ClusterIp)
}

func deleteRuntimeService_old(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	_, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		return
	}
	// here delete network
	multiService, ok := kubeProxyManager.RootMap[runtimeService.Service.Metadata.Name]
	if !ok {
		return
	} else {
		for _, singleService := range multiService {
			err := singleService.deleteSingleService()
			if err != nil {
				fmt.Printf("delete singleService failed: %s\n", err.Error())
				return
			}
		}
	}
	delete(kubeProxyManager.RuntimeServiceMap, runtimeService.Service.Metadata.Name)
}

func updateRuntimeService_old(runtimeService *object.RuntimeService) {
	deleteRuntimeService_old(runtimeService)
	createRuntimeService_old(runtimeService)
}
