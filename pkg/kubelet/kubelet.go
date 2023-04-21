package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_kubelet() {
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	go dealPod(podChan)
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
			// fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarPod.Runtime.Status == config.BOUND_STATUS && tarPod.Runtime.Bind == "TEST" {
				CreatePod(&tarPod)
			}
		}
	}
}
