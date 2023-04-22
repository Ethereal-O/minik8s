package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/counter"
	"minik8s/pkg/util/structure"
)

var waitingRs structure.Cmap

var RSExited = make(chan bool)
var RSToExit = make(chan bool)

func Start_rsController() {
	waitingRs.Init(1)
	handleChan := make(chan string, 20)
	rsChan, rsStop := messging.Watch("/"+config.REPLICASET_TYPE, true)
	//podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	go dealRs(rsChan, handleChan)
	//go dealPod(podChan, handleChan)
	go handle(handleChan)
	fmt.Println("Controller start")

	// Wait until Ctrl-C
	<-RSToExit
	rsStop()
	//podStop()
	RSExited <- true
}

func dealRs(rsChan chan string, handleChan chan string) {
	for {
		select {
		case mes := <-rsChan:
			//fmt.Println("[this]", mes)
			if waitingRs.Put(mes2rsName(mes)) {
				handleChan <- mes
			}
		}
	}
}

//func dealPod(podChan chan string, handleChan chan string) {
//	for {
//		select {
//		case mes := <-podChan:
//			// fmt.Println("[this]", mes)
//			var tarPod object.Pod
//			err := json.Unmarshal([]byte(mes), &tarPod)
//			if err != nil {
//				fmt.Println(err.Error())
//			}
//			if tarPod.Runtime.Belong != "" {
//				res := client.Get_object(tarPod.Runtime.Belong, config.REPLICASET_TYPE)
//				if len(res) != 1 {
//					fmt.Println("Cannot find the Rs which the pod belongs to!")
//					continue
//				}
//				if waitingRs.Put(mes2rsName(res[0])) {
//					handleChan <- res[0]
//				}
//			}
//		}
//	}
//}

func handle(handleChan chan string) {
	for {
		select {
		case mes := <-handleChan:
			fmt.Println("Handle RS")
			waitingRs.Get(mes2rsName(mes))
			var tarRs object.ReplicaSet
			err := json.Unmarshal([]byte(mes), &tarRs)
			if tarRs.Runtime.Status != config.RUNNING_STATUS {
				continue
			}
			if err != nil {
				fmt.Println(err.Error())
			}
			targetNum, actualNum := tarRs.Spec.Replicas, 0
			podList := client.GetActivePods()
			var rspodList []object.Pod
			for _, pod := range podList {
				if pod.Runtime.Belong == tarRs.Metadata.Name {
					actualNum++
					rspodList = append(rspodList, pod)
				}
			}
			fmt.Println(targetNum, actualNum)
			if targetNum > actualNum {
				for i := 0; i < targetNum-actualNum; i++ {
					var newPod object.Pod
					uuid := counter.GetUuid()
					newPod.Runtime.Uuid = uuid
					newPod.Kind = config.POD_TYPE
					newPod.Runtime.Belong = tarRs.Metadata.Name
					newPod.Metadata = tarRs.Spec.Template.Metadata
					newPod.Metadata.Name = tarRs.Metadata.Name + "_" + uuid
					newPod.Spec = tarRs.Spec.Template.Spec
					client.AddPod(newPod)
				}
			} else if targetNum < actualNum {
				for i := targetNum; i < actualNum; i++ {
					client.DeletePod(rspodList[i])
				}
			}
		}
	}
}

func mes2rsName(mes string) string {
	var rsObject object.ReplicaSet
	err := json.Unmarshal([]byte(mes), &rsObject)
	if err != nil {
		fmt.Println(err.Error())
		return "error"
	}
	return rsObject.Metadata.Name
}
