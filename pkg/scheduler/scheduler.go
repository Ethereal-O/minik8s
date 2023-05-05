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

var scheduleQueue structure.Queue

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_scheduler() {
	scheduleQueue.Init()
	var policy SchedulePolicy = ScoringPolicy{}

	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)

	go dealPod(podChan)
	go mainLoop(policy)
	fmt.Println("Scheduler start")

	// Wait until Ctrl-C
	<-ToExit
	podStop()
	Exited <- true
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
		mes := scheduleQueue.Pop()
		if mes == nil {
			// scheduleQueue is empty now
			idle = true
			continue
		} else if mes_str, ok := mes.(string); ok {
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes_str), &tarPod)
			if err != nil {
				// Cannot unmarshal, discard it
				idle = false
				continue
			}
			bound := BindPod(&tarPod, policy)
			if bound {
				// Bind pod success!
				idle = false
				continue
			} else {
				// No optional node now
				idle = true
				scheduleQueue.Push(mes_str)
				fmt.Printf("[Scheduler] No optional node for pod %v\n", tarPod.Metadata.Name)
				continue
			}
		} else {
			// Cannot unmarshal, discard it
			idle = false
			continue
		}
	}
}
