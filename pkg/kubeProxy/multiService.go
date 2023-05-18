package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
)

func createMultiService(runtimeService *object.RuntimeService) {
	multiService := make(map[string]*SingleService)
	for _, port := range runtimeService.Service.Spec.Ports {
		var podsInfo []PodInfo
		for _, pod := range runtimeService.Pods {
			podsInfo = append(podsInfo, PodInfo{
				PodName: pod.Metadata.Name,
				PodIP:   pod.Runtime.PodIp,
				PodPort: port.TargetPort,
			})
		}
		singleService := createSingleService(runtimeService, port, podsInfo)
		err := singleService.initSingleService()
		if err != nil {
			fmt.Printf("init singleService failed: %s\n", err.Error())
			return
		}
		multiService[singleService.Name] = singleService
	}
	kubeProxyManager.RootMap[runtimeService.Service.Metadata.Name] = multiService
	kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService
}

func applyAllMultiService() {
	for _, multiService := range kubeProxyManager.RuntimeServiceMap {
		createMultiService(multiService)
	}
}
