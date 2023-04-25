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
)

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_kubelet() {
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	nodeChan, nodeStop := messging.Watch("/"+config.NODE_TYPE, true)

	autoAddNode()
	fmt.Println("Kubelet start")

	go dealPod(podChan)
	go dealNode(nodeChan)

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
			if tarPod.Runtime.Status == config.BOUND_STATUS {
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
			// fmt.Println("[this]", mes)
			var tarNode object.Node
			err := json.Unmarshal([]byte(mes), &tarNode)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarNode.Runtime.Status == config.CREATED_STATUS {
				err = weave.Expose(tarNode.Runtime.ClusterIp + network.Mask)
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
