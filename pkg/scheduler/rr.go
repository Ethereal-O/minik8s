package scheduler

import (
	"minik8s/pkg/util/counter"
	"time"
)

type RRPolicy struct {
}

func (policy RRPolicy) selectNode() string {
	availNodeList := availNode.GetAll()
	for len(availNodeList) == 0 {
		time.Sleep(100 * time.Millisecond)
	}
	return availNodeList[counter.GetRRPolicy()%len(availNodeList)]
}
