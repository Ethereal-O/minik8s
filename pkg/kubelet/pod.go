package kubelet

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/prometheus/client_golang/prometheus"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

func StartPod(pod *object.Pod) bool {
	var containersIdList []string
	// Step 1: Start pause container
	if !StartPauseContainer(pod) {
		return false
	}

	// Step 2: Start common containers
	for _, myContainer := range pod.Spec.Containers {
		result, ID := StartCommonContainer(pod, &myContainer)
		if !result {
			return false
		}
		containersIdList = append(containersIdList, ID)
	}

	pod.Runtime.Status = config.RUNNING_STATUS
	client.AddPod(*pod)
	go ProbeCycle(pod, containersIdList)
	return true
}

func ProbeCycle(pod *object.Pod, containerIdList []string) {
	ctx := context.Background()
	for {
		time.Sleep(5 * time.Second)

		var containerMemoryPercentageList []float64
		var containerCpuPercentageList []float64

		for _, containerId := range containerIdList {
			containerStats, err := Client.ContainerStats(ctx, containerId, false)
			if err != nil {
				panic(err)
			}
			var stats types.StatsJSON
			dec := json.NewDecoder(containerStats.Body)
			if err := dec.Decode(&stats); err != nil {
				panic(err)
			}
			cpuPercent := calculateCPUPercent(stats)
			memPercent := calculateMemPercent(stats)
			//fmt.Printf("[container:%s] (cpuPercent:%.10f),(memPercent:%.10f)\n",
			//	containerId, cpuPercent, memPercent)
			containerMemoryPercentageList = append(containerMemoryPercentageList, memPercent)
			containerCpuPercentageList = append(containerCpuPercentageList, cpuPercent)
		}

		podAvgMemoryPrecentage := avg(containerMemoryPercentageList)
		podAvgCpuPrecentage := avg(containerCpuPercentageList)
		memoryPrecentage.With(prometheus.Labels{"uuid": pod.Runtime.Uuid, "podName": pod.Metadata.Name}).Set(podAvgMemoryPrecentage)
		cpuPrecentage.With(prometheus.Labels{"uuid": pod.Runtime.Uuid, "podName": pod.Metadata.Name}).Set(podAvgCpuPrecentage)

	}
}
