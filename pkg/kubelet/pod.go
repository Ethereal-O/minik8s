package kubelet

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

func StartPod(pod *object.Pod) bool {
	var containersIdList []string
	// Step 1: Start pause container
	result, ID := StartPauseContainer(pod)
	if !result {
		return false
	}
	containersIdList = append(containersIdList, ID)

	// Step 2: Get pod IP from pause container
	inspection, _ := Client.ContainerInspect(Ctx, ID)
	status, _ := inspectionToContainerRuntime(&inspection)
	pod.Runtime.PodIp = status.IP

	// Step 3: Start common containers
	for _, myContainer := range pod.Spec.Containers {
		result, ID := StartCommonContainer(pod, &myContainer)
		if !result {
			return false
		}
		containersIdList = append(containersIdList, ID)
	}

	pod.Runtime.Status = config.RUNNING_STATUS
	pod.Runtime.Containers = containersIdList
	PodToExit[pod.Runtime.Uuid] = make(chan bool)
	PodExited[pod.Runtime.Uuid] = make(chan bool)
	client.AddPod(*pod)
	go ProbeCycle(pod)
	return true
}

func DeletePod(pod *object.Pod) bool {
	// Step 1: Stop probe cycle
	PodToExit[pod.Runtime.Uuid] <- true
	<-PodExited[pod.Runtime.Uuid]
	delete(PodToExit, pod.Runtime.Uuid)
	delete(PodExited, pod.Runtime.Uuid)

	// Step 2: Delete all containers
	for _, containerId := range pod.Runtime.Containers {
		var stopConfig = 1 * time.Second
		err := Client.ContainerStop(Ctx, containerId, &stopConfig)
		if err != nil {
			fmt.Println(err.Error())
			continue // Ignore errors, may happen when container does not exist
		}
		var removeConfig = RemoveConfig{}
		err = Client.ContainerRemove(Ctx, containerId, removeConfig)
		if err != nil {
			fmt.Println(err.Error())
			continue // Ignore errors, may happen when container does not exist
		}
	}

	// Step 3: If pod needs restart, change its status to CREATED
	if pod.Runtime.NeedRestart {
		pod.Runtime.Status = config.CREATED_STATUS
		// Clear NeedRestart bit, unless the pod will always restart after delete
		pod.Runtime.NeedRestart = false
		client.AddPod(*pod)
	}
	return true
}

func ProbeCycle(pod *object.Pod) {
	for {
		select {
		case <-PodToExit[pod.Runtime.Uuid]:
			PodExited[pod.Runtime.Uuid] <- true
			return
		default:
			time.Sleep(1 * time.Second)
			var containerMemoryPercentageList []float64
			var containerCpuPercentageList []float64

			//The pause container should not be calculated and is supposed to be with no error
			for _, containerId := range pod.Runtime.Containers[1:] {
				inspection, err := Client.ContainerInspect(Ctx, containerId)
				if err != nil {
					// Container does not exist, restart the pod!
					PodException(pod)
					return
				}
				status, err := inspectionToContainerRuntime(&inspection)
				if err != nil {
					panic(err)
				}
				if status.State == StateExited {
					// Container has exited, restart the pod!
					PodException(pod)
					return
				}
				containerMemoryPercentageList = append(containerMemoryPercentageList, status.MemPercent)
				containerCpuPercentageList = append(containerCpuPercentageList, status.CpuPercent)
			}
			podAvgMemoryPrecentage := avg(containerMemoryPercentageList)
			podAvgCpuPrecentage := avg(containerCpuPercentageList)
			memoryPrecentage.With(prometheus.Labels{"uuid": pod.Runtime.Uuid, "podName": pod.Metadata.Name}).Set(podAvgMemoryPrecentage)
			cpuPrecentage.With(prometheus.Labels{"uuid": pod.Runtime.Uuid, "podName": pod.Metadata.Name}).Set(podAvgCpuPrecentage)
		}
	}
}

func PodException(pod *object.Pod) {
	// If the pod belongs to an RS, no need to restart, because RS will do it automatically
	pod.Runtime.NeedRestart = pod.Runtime.Belong == ""
	pod.Runtime.Status = config.EXIT_STATUS
	client.AddPod(*pod)

	// Wait for DeletePod()
	<-PodToExit[pod.Runtime.Uuid]
	PodExited[pod.Runtime.Uuid] <- true
}
