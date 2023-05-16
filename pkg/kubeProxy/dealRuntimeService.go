package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"minik8s/pkg/util/tools"
)

func dealRunningRuntimeService(runtimeService *object.RuntimeService) {
	oldRuntimeService, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		fmt.Printf("creating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		createRuntimeService(runtimeService)
	} else if tools.MD5(*oldRuntimeService) != tools.MD5(*runtimeService) {
		fmt.Printf("updating runtimeService %s\n", runtimeService.Service.Metadata.Name)
		updateRuntimeService(runtimeService)
	} else {
		fmt.Printf("duplicated runtimeService %s\n", runtimeService.Service.Metadata.Name)
	}
}

func dealExitRuntimeService(runtimeService *object.RuntimeService) {
	fmt.Printf("deleting runtimeService %s\n", runtimeService.Service.Metadata.Name)
	deleteRuntimeService(runtimeService)
}

func createRuntimeService(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	if runtimeService.Status == services.SERVICE_STATUS_INIT {
		createDir(services.SERVICE_NGINX_PATH_PREFIX + "/" + runtimeService.Service.Metadata.Name)
	} else if runtimeService.Status == services.SERVICE_STATUS_RUNNING {
		kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService

		applyService(runtimeService)
	}
}

func deleteRuntimeService(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	_, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		return
	}
	delete(kubeProxyManager.RuntimeServiceMap, runtimeService.Service.Metadata.Name)
	deleteDir(services.SERVICE_NGINX_PATH_PREFIX + "/" + runtimeService.Service.Metadata.Name)
}

func updateRuntimeService(runtimeService *object.RuntimeService) {
	//deleteRuntimeService(runtimeService)
	createRuntimeService(runtimeService)
}
