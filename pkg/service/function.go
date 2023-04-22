package service

import "minik8s/pkg/object"

func compareOldAndNewPods(oldPods []*object.Pod, newPods []*object.Pod) bool {
	if len(oldPods) != len(newPods) {
		return true
	}
	for _, oldPod := range oldPods {
		founded := false
		for _, newPod := range newPods {
			if oldPod.Metadata.Name == newPod.Metadata.Name {
				founded = true
				break
			}
		}
		if !founded {
			return true
		}
	}
	return false
}
