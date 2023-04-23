package services

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

func Filter[T any](slice []T, condition func(T) bool) ([]T, []T) {
	var filtered []T
	var differed []T
	for _, item := range slice {
		if condition(item) {
			filtered = append(filtered, item)
		} else {
			differed = append(differed, item)
		}
	}
	return filtered, differed
}

func ForEach[T any](slice []T, action func(T)) {
	for _, item := range slice {
		action(item)
	}
}
