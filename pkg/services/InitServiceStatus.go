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
	go serviceStatus.pollLoop()
}

func (serviceStatus *ServiceStatus) pollLoop() {
	for {
		//if runtimeService.canPollSend {
		//	runtimeService.pollSend()
		//}
		time.Sleep(time.Second * 5)
	}
}

func (serviceStatus *ServiceStatus) poll() {

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
			if v.Metadata.Name == pod.Metadata.Name {

				return true
			}
		}
		return false
	})

	// if no change or has max pods, return
	if len(brokenPods) == 0 && (len(serviceStatus.Pods) >= MAX_PODS || len(serviceStatus.Pods) == len(filterPods)) {
		return
	}

	//// mark normal pods if need to replace net
	//ForEach(normalPods, func(pod object.Pod) {
	//	for _, v := range filterPods {
	//		if v.Metadata.Name == pod.Metadata.Name {
	//			if v.Runtime.Status == config.NEED_REPLACE_STATUS {
	//				pod.Runtime.Status = config.NEED_REPLACE_STATUS
	//			}
	//		}
	//	}
	//})

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

	////更新serviceConfig 并上传
	////如果没有选取到pod, 需要报错，同时如果已经是错误的不需要在去更新etcd
	//var replace []object.PodNameAndIp
	//updateEtcd := false
	//if len(service.pods) == 0 {
	//	if len(service.serviceConfig.Spec.PodNameAndIps) != 0 {
	//		service.serviceConfig.Spec.PodNameAndIps = replace
	//		service.serviceConfig.Status.Phase = object.Failed
	//		service.serviceConfig.Status.Err = NoPodsError
	//		updateEtcd = true
	//	} else {
	//		if isInit {
	//			//第一次就没选到pod
	//			service.serviceConfig.Spec.PodNameAndIps = replace
	//			service.serviceConfig.Status.Phase = object.Failed
	//			service.serviceConfig.Status.Err = NoPodsError
	//		}
	//	}
	//} else {
	//	//比较旧的和新的pod
	//	if compareOldAndNew(oldPods, service.pods) {
	//		//先生成replace
	//		for _, val := range service.pods {
	//			replace = append(replace, object.PodNameAndIp{Name: val.Name, Ip: val.Status.PodIP})
	//		}
	//		//更新serviceConfig
	//		service.serviceConfig.Spec.PodNameAndIps = replace
	//		service.serviceConfig.Status.Phase = object.Running
	//		if service.serviceConfig.Status.Err == NoPodsError {
	//			service.serviceConfig.Status.Err = ""
	//		}
	//		updateEtcd = true
	//	}
	//}
	////更新etcd
	//if isInit {
	//	//第一次是一定要更新的
	//	err = service.Client.UpdateRuntimeService(service.serviceConfig)
	//	return err
	//} else if updateEtcd {
	//	err = service.Client.UpdateRuntimeService(service.serviceConfig)
	//}
	//return err
}
