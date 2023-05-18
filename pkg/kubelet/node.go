package kubelet

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/network"
	"time"
)

func StartNode() {
	ip, err := network.GetHostIp()
	if err != nil {
		fmt.Println("[Kubelet] Cannot obtain host IP!")
		panic(err)
	} else {
		fmt.Printf("[Kubelet] Obtained host IP: %v\n", ip)
	}

	var node object.Node
	node.Kind = "Node"
	node.Metadata.Name = "Node_" + ip
	node.Spec.Ip = ip

	client.AddNode(node)
}

func DeleteNode(node *object.Node) bool {
	// Step 1: Stop probe cycle
	NodeToExit <- true
	<-NodeExited

	// Step 2: Delete all pods
	for _, pod := range client.GetActivePods() {
		if pod.Runtime.Bind == node.Metadata.Name {
			pod.Runtime.Status = config.EXIT_STATUS
			client.AddPod(pod)
		}
	}

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
