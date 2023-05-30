package kubelet

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

func DeleteNode() bool {
	// Step 0: Check if my node has started
	myNode := getMyNode()
	if myNode == nil {
		return true
	}

	// Step 1: Stop probe cycle
	NodeToExit <- true
	<-NodeExited

	// Step 2: Delete all pods
	for _, pod := range client.GetActivePods() {
		if pod.Runtime.Bind == myNode.Metadata.Name {
			pod.Runtime.Status = config.EXIT_STATUS
			client.AddPod(pod)
			<-PodDeleted[pod.Runtime.Uuid] // Wait for pod delete
			delete(PodDeleted, pod.Runtime.Uuid)
		}
	}

	// Step 3: Stop current node
	myNode.Runtime.Status = config.EXIT_STATUS
	client.AddNode(*myNode)
	fmt.Printf("[Kubelet] Node %v deleted!\n", myNode.Metadata.Name)
	return true
}

func NodeProbeCycle(node *object.Node) {
	for {
		select {
		case <-NodeToExit:
			NodeExited <- true
			return
		default:
			// Get idle CPU percentage
			percentIdle, err := cpu.Percent(time.Second, true)
			if err != nil {
				fmt.Println("Failed to retrieve CPU usage:", err)
				return
			}
			var nodeCpu float64 = 0
			for _, percent := range percentIdle {
				nodeCpu += (100 - percent) * 1e7
			}

			// Get free memory
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				fmt.Println("Failed to retrieve memory information:", err)
				return
			}
			var nodeMemory = float64(memInfo.Available)

			nodeAvailableMemory.With(prometheus.Labels{"uuid": node.Runtime.Uuid, "nodeName": node.Metadata.Name}).Set(nodeMemory)
			nodeAvailableCpu.With(prometheus.Labels{"uuid": node.Runtime.Uuid, "nodeName": node.Metadata.Name}).Set(nodeCpu)
		}
	}
}
