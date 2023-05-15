package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/counter"
	"time"
)

var DSControllerExited = make(chan bool)
var DSControllerToExit = make(chan bool)

var DSExited = make(map[string]chan bool)
var DSToExit = make(map[string]chan bool)

func Start_dsController() {
	dsChan, dsStop := messging.Watch("/"+config.DAEMONSET_TYPE, true)
	go dealDs(dsChan)
	fmt.Println("DaemonSet Controller start")

	// Wait until Ctrl-C
	<-DSControllerToExit
	dsStop()
	DSControllerExited <- true
}

func dealDs(dsChan chan string) {
	for {
		select {
		case mes := <-dsChan:
			if mes == "hello" {
				continue
			}
			var tarDs object.DaemonSet
			err := json.Unmarshal([]byte(mes), &tarDs)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarDs.Runtime.Status == config.CREATED_STATUS {
				tarDs.Runtime.Status = config.RUNNING_STATUS
				client.AddDaemonSet(tarDs)
				DSExited[tarDs.Metadata.Name] = make(chan bool)
				DSToExit[tarDs.Metadata.Name] = make(chan bool)
				go DSCycle(tarDs.Metadata.Name)
			} else if tarDs.Runtime.Status == config.EXIT_STATUS {
				DSToExit[tarDs.Metadata.Name] <- true
				<-DSExited[tarDs.Metadata.Name]
				delete(DSToExit, tarDs.Metadata.Name)
				delete(DSExited, tarDs.Metadata.Name)

				dspodList, actualNum := object.GetPodsOfDS(&tarDs, client.GetActivePods())
				for i := 0; i < actualNum; i++ {
					client.DeletePod(dspodList[i])
				}
			}
		}
	}
}

func DSCycle(reName string) {
	for {
		select {
		case <-DSToExit[reName]:
			DSExited[reName] <- true
			return
		default:
			ds := client.GetDaemonSetByKey(reName)[0]
			time.Sleep(1 * time.Second)
			dspodList, _ := object.GetPodsOfDS(&ds, client.GetActivePods())
			nodeList := client.GetActiveNodes()

			// Create pods
			for _, node := range nodeList {
				podExist := false
				for _, pod := range dspodList {
					if pod.Metadata.Name == object.DSPodFullName(&ds, &node) {
						podExist = true
						break
					}
				}
				if podExist {
					continue
				}
				var newPod object.Pod
				uuid := counter.GetUuid()
				newPod.Runtime.Uuid = uuid
				newPod.Kind = config.POD_TYPE
				newPod.Runtime.Status = config.BOUND_STATUS // Bind the pod now to avoid scheduler
				newPod.Runtime.Bind = node.Metadata.Name
				newPod.Runtime.ClusterIp = node.Runtime.ClusterIp
				newPod.Runtime.Belong = ds.Metadata.Name
				newPod.Metadata = ds.Spec.Template.Metadata
				newPod.Metadata.Name = object.DSPodFullName(&ds, &node)
				newPod.Spec = ds.Spec.Template.Spec
				client.AddPod(newPod)
			}

			// Delete pods
			for _, pod := range dspodList {
				nodeExist := false
				for _, node := range nodeList {
					if pod.Metadata.Name == object.DSPodFullName(&ds, &node) {
						nodeExist = true
						break
					}
				}
				if nodeExist {
					continue
				}
				client.DeletePod(pod)
			}
		}
	}
}
