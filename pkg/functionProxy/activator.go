package functionProxy

import (
	"minik8s/pkg/client"
)

func activate(rsName string) {
	tarRs := client.GetReplicaSetByKey(rsName)[0]
	if tarRs.Spec.Replicas == 0 {
		tarRs.Spec.Replicas = 1
		client.AddReplicaSet(tarRs)
	}
}
