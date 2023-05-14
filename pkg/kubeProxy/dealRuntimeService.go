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
	//// here create network
	//fmt.Printf("creating network for runtimeService %s at %s!\n", runtimeService.Service.Metadata.Name, runtimeService.Service.Runtime.ClusterIp)
	//multiService := make(map[string]*SingleService)
	//for _, port := range runtimeService.Service.Spec.Ports {
	//	var podsInfo []PodInfo
	//	for _, pod := range runtimeService.Pods {
	//		podsInfo = append(podsInfo, PodInfo{
	//			PodName: pod.Metadata.Name,
	//			PodIP:   pod.Runtime.PodIp,
	//			PodPort: port.TargetPort,
	//		})
	//	}
	//	singleService := createSingleService(runtimeService, port, podsInfo)
	//	err := singleService.initSingleService()
	//	if err != nil {
	//		fmt.Printf("init singleService failed: %s\n", err.Error())
	//		return
	//	}
	//	multiService[singleService.Name] = singleService
	//}
	//kubeProxyManager.RootMap[runtimeService.Service.Metadata.Name] = multiService
	//kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService
	//fmt.Printf("creating network done for runtimeService %s at %s!\n", runtimeService.Service.Metadata.Name, runtimeService.Service.Runtime.ClusterIp)
	if runtimeService.Status == services.SERVICE_STATUS_INIT {
		createDir(services.SERVICE_NGINX_PATH_PREFIX + "/" + runtimeService.Service.Metadata.Name)
	} else if runtimeService.Status == services.SERVICE_STATUS_RUNNING {
		runtimeService.Status = services.SERVICE_STATUS_RUNNING
		kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService

		updateServiceNginxConfig(runtimeService)
		fmt.Println("write nginx config finished")
		reloadNginxConfig(services.SERVICE_CONTAINER_PREFIX + runtimeService.Service.Metadata.Name)
		fmt.Println("reload nginx config finished")
	}
}

func deleteRuntimeService(runtimeService *object.RuntimeService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	_, ok := kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name]
	if !ok {
		return
	}
	//// here delete network
	//multiService, ok := kubeProxyManager.RootMap[runtimeService.Service.Metadata.Name]
	//if !ok {
	//	return
	//} else {
	//	for _, singleService := range multiService {
	//		err := singleService.deleteSingleService()
	//		if err != nil {
	//			fmt.Printf("delete singleService failed: %s\n", err.Error())
	//			return
	//		}
	//	}
	//}
	delete(kubeProxyManager.RuntimeServiceMap, runtimeService.Service.Metadata.Name)
	deleteDir(services.SERVICE_NGINX_PATH_PREFIX + "/" + runtimeService.Service.Metadata.Name)
}

func updateRuntimeService(runtimeService *object.RuntimeService) {
	//deleteRuntimeService(runtimeService)
	createRuntimeService(runtimeService)
}
