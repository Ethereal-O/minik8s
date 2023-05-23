package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/network"
	"minik8s/pkg/util/weave"
	"time"
)

var Exited = make(chan bool)
var ToExit = make(chan bool)
var PodToExit = make(map[string]chan bool)
var PodExited = make(map[string]chan bool)
var PodDeleted = make(map[string]chan bool)
var NodeToExit = make(chan bool)
var NodeExited = make(chan bool)

func Start_kubelet() {
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	nodeChan, nodeStop := messging.Watch("/"+config.NODE_TYPE, true)

	go dealPod(podChan)
	go dealNode(nodeChan)
	go start_monitor()

	time.Sleep(5 * time.Second)
	StartNode()
	fmt.Println("Kubelet start")

	// Wait until Ctrl-C
	<-ToExit
	DeleteNode()
	podStop()
	nodeStop()
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
			ip, _ := network.GetHostIp()
			if tarPod.Runtime.Status == config.BOUND_STATUS && tarPod.Runtime.Bind == "Node_"+ip {
				go func(pod *object.Pod) {
					started := StartPod(pod)
					if started {
						fmt.Printf("[Kubelet] Pod %v started!\n", pod.Metadata.Name)
					} else {
						fmt.Printf("[Kubelet] Failed to start pod %v!\n", pod.Metadata.Name)
					}
				}(&tarPod)
			} else if tarPod.Runtime.Status == config.EXIT_STATUS && tarPod.Runtime.Bind == "Node_"+ip {
				go func(pod *object.Pod) {
					deleted := DeletePod(pod)
					PodDeleted[pod.Runtime.Uuid] <- true
					if deleted {
						fmt.Printf("[Kubelet] Pod %v deleted!\n", pod.Metadata.Name)
					} else {
						fmt.Printf("[Kubelet] Failed to delete pod %v!\n", pod.Metadata.Name)
					}
				}(&tarPod)
			}
		}
	}
}

func dealNode(nodeChan chan string) {
	for {
		select {
		case mes := <-nodeChan:
			if mes == "hello" {
				continue
			}
			//fmt.Println("[this]", mes)
			var tarNode object.Node
			err := json.Unmarshal([]byte(mes), &tarNode)
			if err != nil {
				fmt.Println(err.Error())
			}
			ip, _ := network.GetHostIp()
			if tarNode.Runtime.Status == config.CREATED_STATUS && tarNode.Metadata.Name == "Node_"+ip {
				err = weave.Expose(tarNode.Runtime.ClusterIp + network.Mask)
				if err != nil {
					fmt.Println("[Kubelet] Failed to start node!")
					fmt.Println(err.Error())
				} else {
					tarNode.Runtime.Status = config.RUNNING_STATUS
					client.AddNode(tarNode)
					fmt.Println("[Kubelet] Node started!")
					go NodeProbeCycle(&tarNode)
				}
			}
		}
	}
}
