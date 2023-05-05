package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

var HpaControllerExited = make(chan bool)
var HpaControllerToExit = make(chan bool)

var HpaExited = make(map[string]chan bool)
var HpaToExit = make(map[string]chan bool)

func Start_hpaController() {
	hpaChan, stopFunc := messging.Watch("/"+config.AUTOSCALER_TYPE, true)
	dealHpa(hpaChan)
	fmt.Println("Autoscaler Controller start")

	// Wait until Ctrl-C
	<-HpaControllerToExit
	stopFunc()
	HpaControllerExited <- true
}

func dealHpa(hpaChan chan string) {
	for {
		select {
		case mes := <-hpaChan:
			if mes == "hello" {
				continue
			}
			// fmt.Println("[this]", mes)
			var tarAutoScaler object.AutoScaler
			json.Unmarshal([]byte(mes), &tarAutoScaler)
			if tarAutoScaler.Runtime.Status == config.CREATED_STATUS {
				tarAutoScaler.Runtime.Status = config.RUNNING_STATUS
				client.AddAutoScaler(tarAutoScaler)
				HpaExited[tarAutoScaler.Metadata.Name] = make(chan bool)
				HpaToExit[tarAutoScaler.Metadata.Name] = make(chan bool)
				go HpaCycle(tarAutoScaler)

			} else if tarAutoScaler.Runtime.Status == config.EXIT_STATUS {
				client.DeleteAutoScaler(tarAutoScaler)
				HpaToExit[tarAutoScaler.Metadata.Name] <- true
				<-HpaExited[tarAutoScaler.Metadata.Name]
				delete(HpaToExit, tarAutoScaler.Metadata.Name)
				delete(HpaExited, tarAutoScaler.Metadata.Name)
			}
		}
	}
}

func HpaCycle(autoScaler object.AutoScaler) {
	var t = autoScaler.Spec.Interval
	ticker := time.NewTicker(time.Duration(t) * time.Second)

	var cpuUnderLimit bool
	var memoryUnderLimit bool
	maxRepilcas := autoScaler.Spec.MaxReplicas
	minReplicas := autoScaler.Spec.MinReplicas

	for {
		select {
		case <-HpaToExit[autoScaler.Metadata.Name]:
			HpaExited[autoScaler.Metadata.Name] <- true
			return
		case <-ticker.C:
			rsList := client.GetReplicaSetByKey(autoScaler.Spec.ScaleTargetRef.Name)
			tarRs := rsList[0]
			fmt.Println("[Hpa ticker!!!!!!!!!!!!!!!!]", tarRs.Metadata.Name)
			rspodList, _ := object.GetPodsOfRS(&tarRs, client.GetActivePods())
			var rspodUuidList []string
			for _, pod := range rspodList {
				rspodUuidList = append(rspodUuidList, pod.Runtime.Uuid)
			}
			fmt.Println("[calculate start]")
			cpuUnderLimit = judge("cpu", autoScaler.Spec.TargetCPUUtilizationStrategy,
				autoScaler.Spec.TargetCPUUtilizationPercentage, rspodUuidList)
			memoryUnderLimit = judge("memory", autoScaler.Spec.TargetMemoryUtilizationStrategy,
				autoScaler.Spec.TargetMemoryUtilizationPercentage, rspodUuidList)
			fmt.Println("[calculate end]")

			if cpuUnderLimit && memoryUnderLimit {
				if tarRs.Spec.Replicas > minReplicas {
					tarRs.Spec.Replicas--
					fmt.Println("[-------------------------]")
					client.AddReplicaSet(tarRs)
				}
			} else {
				if tarRs.Spec.Replicas < maxRepilcas {
					tarRs.Spec.Replicas++
					fmt.Println("[+++++++++++++++++++++++++]")
					client.AddReplicaSet(tarRs)
				}
			}
		}
	}
}
