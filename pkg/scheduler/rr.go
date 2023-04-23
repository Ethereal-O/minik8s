package scheduler

import (
	"minik8s/pkg/util/counter"
)

type RRPolicy struct {
}

func (policy RRPolicy) selectNode() string {
	availNodeList := availNode.GetAll()
	if len(availNodeList) == 0 {
		return ""
	}
	if node, ok := availNodeList[counter.GetRRPolicy()%len(availNodeList)].(string); ok {
		return node
	} else {
		return ""
	}
}
