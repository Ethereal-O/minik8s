package kubeProxy

import (
	"fmt"
	"math/rand"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
)

func dealRunningVirtualService(virtualService *object.VirtualService) {
	oldVirtualService, ok := kubeProxyManager.VirtualServiceMap[virtualService.Metadata.Name]
	if !ok {
		fmt.Printf("creating virtualService %s\n", virtualService.Metadata.Name)
		createVirtualService(virtualService)
	} else if tools.MD5(*oldVirtualService) != tools.MD5(*virtualService) {
		fmt.Printf("updating virtualService %s\n", virtualService.Metadata.Name)
		updateVirtualService(virtualService)
	} else {
		fmt.Printf("duplicated virtualService %s\n", virtualService.Metadata.Name)
	}
}

func dealExitVirtualService(virtualService *object.VirtualService) {
	fmt.Printf("deleting virtualService %s\n", virtualService.Metadata.Name)
	deleteVirtualService(virtualService)
}

func createVirtualService(virtualService *object.VirtualService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	// if no service,return
	runtimeService, ok := kubeProxyManager.RuntimeServiceMap[virtualService.Spec.Service]
	if !ok {
		fmt.Printf("creating virtual service %s failed: no service %s\n", virtualService.Metadata.Name, virtualService.Spec.Service)
		return
	}
	kubeProxyManager.VirtualServiceMap[virtualService.Metadata.Name] = virtualService
	// update matchPodMap
	oldPodMap := kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp]
	for _, selector := range virtualService.Spec.Selector {
		for podName, pod := range oldPodMap {
			flag := true
			for k, v := range selector.MatchLabels {
				podLabel, ok := pod.Pod.Metadata.Labels[k]
				if !ok || !selectorMatch(podLabel, v, virtualService.Spec.Type) {
					flag = false
					break
				}
			}
			if !flag {
				continue
			}
			oldPodMap[podName].PodWeight = selector.Weight
		}
	}
	kubeProxyManager.PodMatchMap[runtimeService.Service.Runtime.ClusterIp] = oldPodMap
}

func deleteVirtualService(virtualService *object.VirtualService) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	_, ok := kubeProxyManager.VirtualServiceMap[virtualService.Metadata.Name]
	if !ok {
		return
	}
	delete(kubeProxyManager.VirtualServiceMap, virtualService.Metadata.Name)
}

func updateVirtualService(virtualService *object.VirtualService) {
	deleteVirtualService(virtualService)
	createVirtualService(virtualService)
}

func getTargetIp(clusterIP string) string {
	// not clusterIP
	pods, ok := kubeProxyManager.PodMatchMap[clusterIP]
	if !ok || len(pods) == 0 {
		return clusterIP
	}

	sumWeight := 0
	for _, pod := range pods {
		sumWeight += pod.PodWeight
	}

	// random select
	if sumWeight == 0 {
		idx := rand.Intn(len(pods))
		var selectedIp string
		for _, pod := range pods {
			if idx == 0 {
				selectedIp = pod.Pod.Runtime.PodIp
				break
			}
			idx--
		}
		fmt.Printf("selecting ip %v for service %v\n", selectedIp, clusterIP)
		return selectedIp
	}

	// weight based select
	numWeight := rand.Intn(sumWeight) + 1
	sumWeight = 0
	for _, pod := range pods {
		sumWeight += pod.PodWeight
		if sumWeight >= numWeight {
			fmt.Printf("selecting ip %v for service %v\n", pod.Pod.Runtime.PodIp, clusterIP)
			return pod.Pod.Runtime.PodIp
		}
	}

	fmt.Printf("find enpoint by weight fail, clusterIP:%v\n", clusterIP)

	return clusterIP
}
