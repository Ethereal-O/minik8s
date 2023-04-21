package scheduler

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func BindPod(pod *object.Pod) {
	pod.Runtime.Status = config.BOUND_STATUS
	pod.Runtime.Bind = "TEST"
	client.AddPod(*pod)
}
