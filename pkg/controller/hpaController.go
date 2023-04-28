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

var HpaExited = make(chan bool)
var HpaToExit = make(chan bool)

var DealingExited = make(map[string]chan bool)
var DealingToExit = make(map[string]chan bool)

func Start_hpaController() {
	hpaChan, stopFunc := messging.Watch("/"+config.AUTOSCALER_TYPE, true)
	dealHpa(hpaChan)
	fmt.Println("Autoscaler Controller start")

	// Wait until Ctrl-C
	<-HpaToExit
	stopFunc()
	HpaExited <- true
}

func dealHpa(hpaChan chan string) {
	for {
		select {
		case mes := <-hpaChan:
			if mes == "hello" {
				continue
			}
			fmt.Println("[this]", mes)
			var tarAutoScaler object.AutoScaler
			json.Unmarshal([]byte(mes), &tarAutoScaler)
			if tarAutoScaler.Runtime.Status == config.CREATED_STATUS {
				tarAutoScaler.Runtime.Status = config.RUNNING_STATUS
				client.AddAutoScaler(tarAutoScaler)
				DealingExited[tarAutoScaler.Metadata.Name] = make(chan bool)
				DealingToExit[tarAutoScaler.Metadata.Name] = make(chan bool)
				go dealingCycle(tarAutoScaler)

			} else if tarAutoScaler.Runtime.Status == config.EXIT_STATUS {
				client.DeleteAutoScaler(tarAutoScaler)
				DealingToExit[tarAutoScaler.Metadata.Name] <- true
				<-DealingExited[tarAutoScaler.Metadata.Name]
				delete(DealingToExit, tarAutoScaler.Metadata.Name)
				delete(DealingExited, tarAutoScaler.Metadata.Name)
			}
		}
	}
}

func dealingCycle(autoScaler object.AutoScaler) {
	var t = autoScaler.Spec.Interval
	ticker := time.NewTicker(time.Duration(t) * time.Second)

	var cpuUnderLimit bool
	var memoryUnderLimit bool
	maxRepilcas := autoScaler.Spec.MaxReplicas
	minReplicas := autoScaler.Spec.MinReplicas

	for {
		select {
		case <-DealingToExit[autoScaler.Metadata.Name]:
			DealingExited[autoScaler.Metadata.Name] <- true
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
