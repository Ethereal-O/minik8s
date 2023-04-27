package services

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/tools"
	"time"
)

func InitRuntimeService(runtimeService *object.RuntimeService) {
	selectPods(runtimeService)
	startPoll(runtimeService)
}

func startPoll(runtimeService *object.RuntimeService) {
	runtimeService.Timer = *time.NewTicker(CHECK_PODS_TIME_INTERVAL)
	go pollLoop(runtimeService)
}

func pollLoop(runtimeService *object.RuntimeService) {
	defer runtimeService.Timer.Stop()
	for {
		select {
		case <-runtimeService.Timer.C:
			poll(runtimeService)
		}
	}
}

func poll(runtimeService *object.RuntimeService) {
	runtimeService.Lock.Lock()
	defer runtimeService.Lock.Unlock()
	selectPods(runtimeService)
}

func selectPods(runtimeService *object.RuntimeService) {
	// get all pods and selector
	selector := runtimeService.Service.Spec.Selector
	allPods := client.GetAllPods()

	// apply filter to get new pods
	filterPods, _ := tools.Filter(allPods, func(pod object.Pod) bool {
		if pod.Runtime.Status != config.RUNNING_STATUS {
			return false
		}
		for k, v := range selector {
			podLabel, ok := pod.Metadata.Labels[k]
			if !ok || v != podLabel {
				return false
			}
		}
		return true
	})

	// apply filter to get broken pods
	normalPods, brokenPods := tools.Filter(runtimeService.Pods, func(pod object.Pod) bool {
		for _, v := range filterPods {
			if v.Metadata.Name == pod.Metadata.Name && v.Runtime.ClusterIp == pod.Runtime.ClusterIp {
				return true
			}
		}
		return false
	})

	// if no change or has max pods, return
	if len(brokenPods) == 0 && (len(runtimeService.Pods) >= MAX_PODS || len(runtimeService.Pods) == len(filterPods)) {
		return
	}

	// try to fill the pods to max_pods
	_, differPods := tools.Filter(filterPods, func(pod object.Pod) bool {
		for _, v := range normalPods {
			if v.Metadata.Name == pod.Metadata.Name {
				return true
			}
		}
		return false
	})

	// fill normalPods to max pods
	for _, pod := range differPods {
		if len(normalPods) >= MAX_PODS {
			break
		}
		normalPods = append(normalPods, pod)
	}

	// update service status
	runtimeService.Pods = normalPods
	fmt.Printf("service %s update pods num: %d\n", runtimeService.Service.Metadata.Name, len(runtimeService.Pods))

	// update service config
	client.AddRuntimeService(*runtimeService)
}
