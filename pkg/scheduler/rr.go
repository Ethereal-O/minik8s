package scheduler

import (
	"minik8s/pkg/object"
	"minik8s/pkg/util/counter"
)

type RRPolicy struct {
}

func (policy RRPolicy) selectNode(pod *object.Pod, nodes []*object.Node) string {
	// No optional nodes
	if len(nodes) == 0 {
		return ""
	}
	return nodes[counter.GetRRPolicy()%len(nodes)].Metadata.Name
}
