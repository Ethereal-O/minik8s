package kubelet

import (
	"minik8s/pkg/object"
)

func CreatePod(pod *object.Pod) {
	// Step 1: Start pause container
	StartPauseContainer(pod)

	// Step 2: Start common containers
	for _, myContainer := range pod.Spec.Containers {
		StartCommonContainer(pod, &myContainer)
	}
}
