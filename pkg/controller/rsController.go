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

var RSControllerExited = make(chan bool)
var RSControllerToExit = make(chan bool)

var RSExited = make(map[string]chan bool)
var RSToExit = make(map[string]chan bool)

func Start_rsController() {
	rsChan, rsStop := messging.Watch("/"+config.REPLICASET_TYPE, true)
	go dealRs(rsChan)
	fmt.Println("Replicaset Controller start")

	// Wait until Ctrl-C
	<-RSControllerToExit
	rsStop()
	RSControllerExited <- true
}

func dealRs(rsChan chan string) {
	for {
		select {
		case mes := <-rsChan:
			if mes == "hello" {
				continue
			}
			var tarRs object.ReplicaSet
			err := json.Unmarshal([]byte(mes), &tarRs)
			if err != nil {
				fmt.Println(err.Error())
			}
			if tarRs.Runtime.Status == config.CREATED_STATUS {
				tarRs.Runtime.Status = config.RUNNING_STATUS
				client.AddReplicaSet(tarRs)
				RSExited[tarRs.Metadata.Name] = make(chan bool)
				RSToExit[tarRs.Metadata.Name] = make(chan bool)
				go RSCycle(tarRs.Metadata.Name)
			} else if tarRs.Runtime.Status == config.EXIT_STATUS {
				RSToExit[tarRs.Metadata.Name] <- true
				<-RSExited[tarRs.Metadata.Name]
				delete(RSToExit, tarRs.Metadata.Name)
				delete(RSExited, tarRs.Metadata.Name)

				rspodList, actualNum := object.GetPodsOfRS(&tarRs, client.GetActivePods())
				for i := 0; i < actualNum; i++ {
					client.DeletePod(rspodList[i])
				}
			}
		}
	}
}

func RSCycle(reName string) {
	for {
		select {
		case <-RSToExit[reName]:
			RSExited[reName] <- true
			return
		default:
			rs := client.GetReplicaSetByKey(reName)[0]
			time.Sleep(1 * time.Second)
			targetNum := rs.Spec.Replicas
			rspodList, actualNum := object.GetPodsOfRS(&rs, client.GetActivePods())
			if targetNum > actualNum {
				for i := 0; i < targetNum-actualNum; i++ {
					var newPod object.Pod
					uuid := counter.GetUuid()
					newPod.Runtime.Uuid = uuid
					newPod.Kind = config.POD_TYPE
					newPod.Runtime.Belong = rs.Metadata.Name
					newPod.Metadata = rs.Spec.Template.Metadata
					newPod.Metadata.Name = object.RSPodFullName(&rs, &newPod)
					newPod.Spec = rs.Spec.Template.Spec
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
