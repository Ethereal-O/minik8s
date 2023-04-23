package scheduler

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func BindPod(pod *object.Pod, policy SchedulePolicy) bool {
	pod.Runtime.Bind = policy.selectNode()
	if pod.Runtime.Bind != "" {
		pod.Runtime.Status = config.BOUND_STATUS
		client.AddPod(*pod)
		return true
	} else {
		return false
	}
}

type SchedulePolicy interface {
	selectNode() string
}
