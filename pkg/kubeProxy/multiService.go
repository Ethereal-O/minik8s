package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func createMultiService(runtimeService *object.RuntimeService) {
	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES {
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

	}
	if config.SERVICE_POLICY == config.SERVICE_POLICY_MICROSERVICE {
		oldPods, ok := kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp]
		if !ok {
			kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp] = make(map[string]*PodMatch)
			oldPods = kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp]
		}
		for _, pod := range runtimeService.Pods {
			oldPod, ok := oldPods[pod.Metadata.Name]
			if !ok {
				kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp][pod.Metadata.Name] = &PodMatch{
					Pod:       &pod,
					PodWeight: 0,
				}
			} else {
				kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp][pod.Metadata.Name] = &PodMatch{
					Pod:       &pod,
					PodWeight: oldPod.PodWeight,
				}
			}
		}
	}
	kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = runtimeService
}

func applyAllMultiService() {
	for _, multiService := range kubeProxyManager.RuntimeServiceMap {
		createMultiService(multiService)
	}
}
