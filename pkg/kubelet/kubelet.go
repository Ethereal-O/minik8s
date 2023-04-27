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

func Start_kubelet() {
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	nodeChan, nodeStop := messging.Watch("/"+config.NODE_TYPE, true)

	go dealPod(podChan)
	go dealNode(nodeChan)
	go start_monitor()
	
	time.Sleep(5 * time.Second)
	autoAddNode()
	fmt.Println("Kubelet start")

	// Wait until Ctrl-C
	<-ToExit
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
				started := StartPod(&tarPod)
				if started {
					fmt.Printf("[Kubelet] Pod %v started!\n", tarPod.Metadata.Name)
				} else {
					fmt.Printf("[Kubelet] Failed to start pod %v!\n", tarPod.Metadata.Name)
				}
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
			fmt.Println("[this]", mes)
			var tarNode object.Node
			err := json.Unmarshal([]byte(mes), &tarNode)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarNode.Runtime.Status == config.CREATED_STATUS {
				fmt.Println("[weave.Expose start]")
				err = weave.Expose(tarNode.Runtime.ClusterIp + network.Mask)
				fmt.Println("[weave.Expose end]")
				if err != nil {
					fmt.Println("[Kubelet] Failed to start node!")
					fmt.Println(err.Error())
				} else {
					tarNode.Runtime.Status = config.RUNNING_STATUS
					inf, _ := json.Marshal(&tarNode)
					client.Put_object(tarNode.Metadata.Name, string(inf), "Node")
					fmt.Println("[Kubelet] Node started!")
				}
			}
		}
	}
}

func autoAddNode() {
	ip, err := network.GetHostIp()
	if err != nil {
		fmt.Println("[Kubelet] Cannot obtain host IP!")
		panic(err)
	} else {
		fmt.Printf("[Kubelet] Obtained host IP: %v\n", ip)
	}

	var node object.Node
	node.Kind = "Node"
	node.Metadata.Name = "Node_" + ip
	node.Spec.Ip = ip
	inf, _ := json.Marshal(&node)
	client.Put_object(node.Metadata.Name, string(inf), "Node")
}
