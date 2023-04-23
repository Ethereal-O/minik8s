package kubelet

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func StartPod(pod *object.Pod) bool {
	// Step 1: Start pause container
	if !StartPauseContainer(pod) {
		return false
	}

	// Step 2: Start common containers
	for _, myContainer := range pod.Spec.Containers {
		if !StartCommonContainer(pod, &myContainer) {
			return false
		}
	}

	pod.Runtime.Status = config.RUNNING_STATUS
	client.AddPod(*pod)
	return true
}
