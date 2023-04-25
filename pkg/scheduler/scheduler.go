package scheduler

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/structure"
	"time"
)

var availNode structure.Set
var scheduleQueue structure.Queue

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_scheduler() {
	availNode.Init()
	scheduleQueue.Init()
	var policy SchedulePolicy = RRPolicy{}

	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	nodeChan, nodeStop := messging.Watch("/"+config.NODE_TYPE, true)

	go dealPod(podChan)
	go dealNode(nodeChan)
	go mainLoop(policy)
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
			if mes == "hello" {
				continue
			}
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

func dealPod(podChan chan string) {
	for {
		select {
		case mes := <-podChan:
			if mes == "hello" {
				continue
			}
			// fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarPod.Runtime.Status == config.CREATED_STATUS {
				scheduleQueue.Push(mes)
			}
		}
	}
}

func mainLoop(policy SchedulePolicy) {
	idle := false
	for {
		if idle {
			time.Sleep(1 * time.Second)
		}
		mes := scheduleQueue.Front()
		if mes == nil {
			// scheduleQueue is empty now
			idle = true
			continue
		} else if mes_str, ok := mes.(string); ok {
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes_str), &tarPod)
			if err != nil {
				// Cannot unmarshal, discard it
				scheduleQueue.Pop()
				idle = false
				continue
			}
			bound := BindPod(&tarPod, policy)
			if bound {
				// Bind pod success!
				scheduleQueue.Pop()
				idle = false
				continue
			} else {
				// No available node now
				idle = true
				continue
			}
		} else {
			// Cannot unmarshal, discard it
			scheduleQueue.Pop()
			idle = false
			continue
		}
	}
}
