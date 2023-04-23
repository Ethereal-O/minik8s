package services

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

func (serviceStatus *ServiceStatus) InitServiceStatus() {
	serviceStatus.selectPods()
	serviceStatus.startPoll()
}

func (serviceStatus *ServiceStatus) startPoll() {
	serviceStatus.Timer = *time.NewTicker(CHECK_PODS_TIME_INTERVAL)
	go serviceStatus.pollLoop()
}

func (serviceStatus *ServiceStatus) pollLoop() {
	defer serviceStatus.Timer.Stop()
	for {
		select {
		case <-serviceStatus.Timer.C:
			serviceStatus.poll()
		}
	}
}

func (serviceStatus *ServiceStatus) poll() {
	serviceStatus.Lock.Lock()
	defer serviceStatus.Lock.Unlock()
	serviceStatus.selectPods()
}

func (serviceStatus *ServiceStatus) selectPods() {
	// get all pods and selector
	selector := serviceStatus.Service.Spec.Selector
	allPods := client.GetAllPods()

	// apply filter to get new pods
	filterPods, _ := Filter(allPods, func(pod object.Pod) bool {
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
	normalPods, brokenPods := Filter(serviceStatus.Pods, func(pod object.Pod) bool {
		for _, v := range filterPods {
			if v.Metadata.Name == pod.Metadata.Name && v.Runtime.ClusterIp == pod.Runtime.ClusterIp {
				return true
			}
		}
		return false
	})

	// if no change or has max pods, return
	if len(brokenPods) == 0 && (len(serviceStatus.Pods) >= MAX_PODS || len(serviceStatus.Pods) == len(filterPods)) {
		return
	}

	// try to fill the pods to max_pods
	_, differPods := Filter(filterPods, func(pod object.Pod) bool {
		for _, v := range normalPods {
			if v.Metadata.Name == pod.Metadata.Name {
				return false
			}
		}
		return true
	})

	// fill normalPods to max pods
	for _, pod := range differPods {
		if len(normalPods) >= MAX_PODS {
			break
		}
		normalPods = append(normalPods, pod)
	}

	// update service status
	serviceStatus.Pods = normalPods

	// update service config
	client.AddServiceStatus(*serviceStatus)
}
