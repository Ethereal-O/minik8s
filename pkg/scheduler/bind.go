package scheduler

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func BindPod(pod *object.Pod, policy SchedulePolicy) {
	pod.Runtime.Status = config.BOUND_STATUS
	pod.Runtime.Bind = policy.selectNode()
	client.AddPod(*pod)
}

type SchedulePolicy interface {
	selectNode() string
}
