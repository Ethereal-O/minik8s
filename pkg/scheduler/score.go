package scheduler

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
)

type ScoringPolicy struct {
}

func (policy ScoringPolicy) selectNode(pod *object.Pod, nodes []object.Node) string {
	// No optional nodes
	if len(nodes) == 0 {
		return ""
	}

	max_score := 0.0
	max_score_node := ""

	for _, node := range nodes {
		cpu := prometheus_query("cpu", node.Runtime.Uuid)
		mem := prometheus_query("memory", node.Runtime.Uuid)
		cpu_score := float64(cpu / 1e9) * 10.0       // 1 CPU = 10 score
		mem_score := float64(mem / (1024 * 1024 * 1024)) * 5.0 // 1 GB Memory = 5 score

		rs_scale := 1.5
		// pod belongs to a ReplicaSet
		if pod.Runtime.Belong != "" {
			pods := client.GetActivePods()
			for _, existing_pod := range pods {
				// An existing pod is bound to the node and is from the same ReplicaSet
				if existing_pod.Runtime.Belong == pod.Runtime.Belong && existing_pod.Runtime.Bind == node.Metadata.Name {
					rs_scale = 1.0 // ReplicaSet should have pods on various nodes
				}
			}
		}

		score := (cpu_score + mem_score) * rs_scale
		if score > max_score {
			max_score = score
			max_score_node = node.Metadata.Name
		}
	}

	fmt.Printf("[Scheduler] Select Node %v, score = %v\n", max_score_node, max_score)
	return max_score_node
}
