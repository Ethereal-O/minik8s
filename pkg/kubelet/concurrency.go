package kubelet

import (
	"minik8s/pkg/object"
	"sync"
)

type StartResult struct {
	Result bool
	ID     string
	Index  int
}

func StartContainersConcurrently(pod *object.Pod, hostMode bool) (bool, []string) {
	var wg sync.WaitGroup
	containersIdList := make([]string, len(pod.Spec.Containers))
	resultChan := make(chan StartResult, len(pod.Spec.Containers))

	for i, myContainer := range pod.Spec.Containers {
		wg.Add(1)
		go func(container object.Container, index int) {
			defer wg.Done()
			var result bool
			var ID string
			if hostMode {
				result, ID = StartHostContainer(pod, &container)
			} else {
				result, ID = StartCommonContainer(pod, &container)
			}
			resultChan <- StartResult{Result: result, ID: ID, Index: index}
		}(myContainer, i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if !result.Result {
			return false, nil
		}
		containersIdList[result.Index] = result.ID
	}

	return true, containersIdList
}
