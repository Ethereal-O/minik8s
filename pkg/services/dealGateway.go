package services

import (
	"minik8s/pkg/client"
	"minik8s/pkg/object"
)

func dealRunningGateway(gateway *object.Gateway) {
	createGateway(gateway)
}

func dealExitGateway(gateway *object.Gateway) {
	deleteGateway(gateway)
}

func createGateway(gateway *object.Gateway) {
	client.AddReplicaSet(GetGateWayReplicaSet(gateway.MetaData.Name))
	client.AddService(GetGateWayService(gateway.MetaData.Name))
	dnsManager.Lock.Lock()
	defer dnsManager.Lock.Unlock()
	gatewayStatus := object.RuntimeGateway{
		Gateway: *gateway,
		Status:  GATEWAY_STATUS_INIT,
	}
	dnsManager.GatewayMap[gateway.MetaData.Name] = gatewayStatus
}

func deleteGateway(gateway *object.Gateway) {
	serviceList := client.GetServiceByKey(GATEWAY_SERVICE_PREFIX + gateway.MetaData.Name)
	if len(serviceList) == 0 {
		return
	}
	client.DeleteService(serviceList[0])

	replicaSetList := client.GetReplicaSetByKey(GATEWAY_REPLICASET_PREFIX + gateway.MetaData.Name)
	if len(replicaSetList) == 0 {
		return
	}
	client.DeleteReplicaSet(replicaSetList[0])
}
