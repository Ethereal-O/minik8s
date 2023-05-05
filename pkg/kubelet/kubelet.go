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
	"runtime"
	"syscall"
	"time"
)

var Exited = make(chan bool)
var ToExit = make(chan bool)
var PodToExit = make(map[string]chan bool)
var PodExited = make(map[string]chan bool)

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
			node := client.GetNode(ip)
			if tarPod.Runtime.Status == config.BOUND_STATUS && tarPod.Runtime.Bind == node.Metadata.Name {
				started := StartPod(&tarPod, &node)
				if started {
					fmt.Printf("[Kubelet] Pod %v started!\n", tarPod.Metadata.Name)
				} else {
					fmt.Printf("[Kubelet] Failed to start pod %v!\n", tarPod.Metadata.Name)
				}
			} else if tarPod.Runtime.Status == config.EXIT_STATUS && tarPod.Runtime.Bind == node.Metadata.Name {
				deleted := DeletePod(&tarPod, &node)
				if deleted {
					fmt.Printf("[Kubelet] Pod %v deleted!\n", tarPod.Metadata.Name)
				} else {
					fmt.Printf("[Kubelet] Failed to delete pod %v!\n", tarPod.Metadata.Name)
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
			//fmt.Println("[this]", mes)
			var tarNode object.Node
			err := json.Unmarshal([]byte(mes), &tarNode)
			if err != nil {
				fmt.Println(err.Error())
			}
			ip, _ := network.GetHostIp()
			node := client.GetNode(ip)
			if tarNode.Runtime.Status == config.CREATED_STATUS && tarNode.Metadata.Name == node.Metadata.Name {
				err = weave.Expose(tarNode.Runtime.ClusterIp + network.Mask)
				if err != nil {
					fmt.Println("[Kubelet] Failed to start node!")
					fmt.Println(err.Error())
				} else {
					tarNode.Runtime.Status = config.RUNNING_STATUS
					client.AddNode(tarNode)
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
	var cpu = int64(runtime.NumCPU()) * 1e9 / 100 * NodeResourceUsage // NanoCPU
	node.Runtime.Available.Cpu = cpu
	node.Spec.Capacity.Cpu = cpu

	// Sysinfo is only for linux!!
	sysInfo := new(syscall.Sysinfo_t)
	err = syscall.Sysinfo(sysInfo)
	if err != nil {
		fmt.Println("[Kubelet] Cannot obtain host Memory!")
		panic(err)
	}
	var mem = int64(sysInfo.Totalram) / 100 * NodeResourceUsage // Bytes
	node.Runtime.Available.Memory = mem
	node.Spec.Capacity.Memory = mem

	client.AddNode(node)
}
