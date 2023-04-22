package scheduler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/structure"
)

var availNode structure.Set

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_scheduler() {
	availNode.Init()
	var policy SchedulePolicy = RRPolicy{}
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	nodeChan, nodeStop := messging.Watch("/"+config.NODE_TYPE, true)
	go dealPod(podChan, policy)
	go dealNode(nodeChan)
	fmt.Println("Scheduler start")

	// Wait until Ctrl-C
	<-ToExit
	podStop()
	nodeStop()
	Exited <- true
}

func dealNode(nodeChan chan string) {
	for {
		select {
		case mes := <-nodeChan:
			// fmt.Println("[this]", mes)
			var tarNode object.Node
			err := json.Unmarshal([]byte(mes), &tarNode)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarNode.Runtime.Status == config.RUNNING_STATUS {
				availNode.Put(tarNode.Metadata.Name)
			} else if tarNode.Runtime.Status == config.EXIT_STATUS {
				availNode.Del(tarNode.Metadata.Name)
			}
		}
	}
}

func dealPod(podChan chan string, policy SchedulePolicy) {
	for {
		select {
		case mes := <-podChan:
			fmt.Println("scheduler: pod!!")
			// fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarPod.Runtime.Status == config.CREATED_STATUS {
				fmt.Println("scheduler: created pod!!")
				BindPod(&tarPod, policy)
			}
		}
	}
}
