package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/structure"
	"time"
)

var startQueue structure.Queue

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_kubelet() {
	startQueue.Init()

	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)

	go dealPod(podChan)
	go mainLoop()
	fmt.Println("Kubelet start")

	// Wait until Ctrl-C
	<-ToExit
	podStop()
	Exited <- true
}

func dealPod(podChan chan string) {
	for {
		select {
		case mes := <-podChan:
			if mes=="hello" {
				continue
			}
			// fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarPod.Runtime.Status == config.BOUND_STATUS {
				fmt.Println("kubelet: bound pod!!")
				startQueue.Push(mes)
			}
		}
	}
}

func mainLoop() {
	idle := false
	for {
		if idle {
			time.Sleep(1 * time.Second)
		}
		mes := startQueue.Front()
		if mes == nil {
			// startQueue is empty now
			idle = true
			continue
		} else if mes_str, ok := mes.(string); ok {
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes_str), &tarPod)
			if err != nil {
				// Cannot unmarshal, discard it
				startQueue.Pop()
				idle = false
				continue
			}

			started := StartPod(&tarPod)

			if started {
				// Start pod success!
				startQueue.Pop()
				idle = false
				continue
			} else {
				// Invalid pod, skip it
				startQueue.Pop()
				idle = false
				continue
			}
		} else {
			// Cannot unmarshal, discard it
			startQueue.Pop()
			idle = false
			continue
		}
	}
}
