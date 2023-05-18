package functionProxy

import (
	"minik8s/pkg/client"
	"minik8s/pkg/util/structure"
	"strings"
	"time"
)

const (
	scale_boundary       = 5
	flow_reduce_interval = 5
	flow_check_interval  = 30
)

var SvcFlowMap structure.CountMap
var Inf = make(chan string, 10)

func addFlow(svcName string) {
	if !SvcFlowMap.Exist(svcName) {
		Inf <- svcName
		//This is important because after the kubeProxy renew the pod, it needs time to renew service
		time.Sleep(10 * time.Second)
	}
	SvcFlowMap.Add(svcName)
}

func FlowControl() {
	SvcFlowMap.Init()
	for {
		select {
		case mes := <-Inf:
			go svcFlowLoop(mes)
		}
	}
}

func svcFlowLoop(svcName string) {
	minus_ticker := time.NewTicker(time.Duration(flow_reduce_interval) * time.Second)
	check_ticker := time.NewTicker(time.Duration(flow_check_interval) * time.Second)
	preFlow := 0
	tmpFlow := 0
	lastFlag := true
	for {
		select {
		case <-minus_ticker.C:
			SvcFlowMap.Minus(svcName)
		case <-check_ticker.C:
			tmpFlow = SvcFlowMap.Get(svcName)
			if tmpFlow-preFlow >= -scale_boundary && tmpFlow-preFlow <= scale_boundary {
				continue
			} else if tmpFlow-preFlow < -scale_boundary {
				preFlow = tmpFlow
				tarRsName := strings.ReplaceAll(svcName, "service", "rs")
				tarRs := client.GetReplicaSetByKey(tarRsName)[0]
				if tarRs.Spec.Replicas > 1 {
					tarRs.Spec.Replicas--
					client.AddReplicaSet(tarRs)
				} else if tarRs.Spec.Replicas == 1 && lastFlag {
					lastFlag = false
				} else if tarRs.Spec.Replicas == 1 && !lastFlag {
					tarRs.Spec.Replicas--
					client.AddReplicaSet(tarRs)
					SvcFlowMap.Delete(svcName)
					return
				} else {
					SvcFlowMap.Delete(svcName)
					return
				}
			} else if tmpFlow-preFlow > scale_boundary {
				preFlow = tmpFlow
				tarRsName := strings.ReplaceAll(svcName, "service", "rs")
				tarRs := client.GetReplicaSetByKey(tarRsName)[0]
				if tarRs.Spec.Replicas < scale_boundary {
					tarRs.Spec.Replicas++
					client.AddReplicaSet(tarRs)
				}
			}
		}
	}
}
