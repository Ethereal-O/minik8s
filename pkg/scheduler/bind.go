package scheduler

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/resource"
)

func BindPod(pod *object.Pod, policy SchedulePolicy) bool {
	optional_nodes := getOptionalNodes(pod)
	pod.Runtime.Bind = policy.selectNode(pod, optional_nodes)
	if pod.Runtime.Bind != "" {
		pod.Runtime.Status = config.BOUND_STATUS
		client.AddPod(*pod)
		return true
	} else {
		return false
	}
}

type SchedulePolicy interface {
	selectNode(pod *object.Pod, nodes []object.Node) string
}

func getOptionalNodes(pod *object.Pod) []object.Node {
	nodes := client.GetActiveNodes()

	var optional_nodes []object.Node
	for _, node := range nodes {
		cpu := node.Runtime.Available.Cpu
		mem := node.Runtime.Available.Memory
		for _, container := range pod.Spec.Containers {
			cpu -= resource.ConvertCpuToBytes(container.Limits.Cpu)
			mem -= resource.ConvertMemoryToBytes(container.Limits.Memory)
		}
		if cpu > 0 && mem > 0 {
			optional_nodes = append(optional_nodes, node)
		}
	}

	return optional_nodes
}
