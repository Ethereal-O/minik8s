package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/structure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var waitingRs structure.Cmap

func Start_rsController() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	waitingRs.Init(1)
	handleChan := make(chan string, 20)
	rsChan, rsStop := messging.Watch("/"+config.REPLICASET_TYPE, true)
	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	go dealRs(rsChan, handleChan)
	go dealPod(podChan, handleChan)
	go handle(handleChan)

	<-c
	rsStop()
	podStop()
	time.Sleep(2 * time.Second)
	return
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

func dealPod(podChan chan string, handleChan chan string) {
	for {
		select {
		case mes := <-podChan:
			//fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarPod.Belong != "" {
				res := client.Get_object(tarPod.Belong, config.REPLICASET_TYPE)
				if len(res) != 1 {
					fmt.Println("Cannot find the Rs which the pod belongs to!")
					continue
				}
				if waitingRs.Put(mes2rsName(res[0])) {
					handleChan <- res[0]
				}
			}
		}
	}
}

func handle(handleChan chan string) {
	for {
		select {
		case mes := <-handleChan:
			//fmt.Println("[this]", mes)
			waitingRs.Get(mes2rsName(mes))
			var tarRs object.ReplicaSet
			err := json.Unmarshal([]byte(mes), &tarRs)
			if tarRs.Metadata.Status != config.RUNNING_STATUS {
				continue
			}
			if err != nil {
				fmt.Println(err.Error())
			}
			targetNum, actualNum := tarRs.Spec.Replicas, 0
			podList := client.GetRunningPods()
			var rspodList []object.Pod
			for _, pod := range podList {
				if pod.Belong == tarRs.Metadata.Name {
					actualNum++
					rspodList = append(rspodList, pod)
				}
			}
			if targetNum > actualNum {
				for i := 0; i < targetNum-actualNum; i++ {
					var newPod object.Pod
					newPod.Kind = config.POD_TYPE
					newPod.Belong = tarRs.Metadata.Name
					newPod.Metadata = tarRs.Spec.Template.Metadata
					newPod.Metadata.Name = tarRs.Metadata.Name
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
