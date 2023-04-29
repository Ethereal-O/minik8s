package services

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
)

func dealRunningGateway(gateway *object.Gateway) {
	oldGateway, ok := dnsManager.GatewayMap[gateway.Metadata.Name]
	if !ok {
		fmt.Printf("creating gateway %s\n", gateway.Metadata.Name)
		createGateway(gateway)
	} else if tools.MD5(oldGateway.Gateway) != tools.MD5(*gateway) {
		fmt.Printf("updating gateway %s\n", gateway.Metadata.Name)
		updateGateway(gateway)
	} else {
		fmt.Printf("duplicated gateway %s\n", gateway.Metadata.Name)
	}
}

func dealExitGateway(gateway *object.Gateway) {
	deleteGateway(gateway)
}

func createGateway(gateway *object.Gateway) {
	client.AddReplicaSet(GetGateWayReplicaSet(gateway.Metadata.Name))
	client.AddService(GetGateWayService(gateway.Metadata.Name))
	dnsManager.Lock.Lock()
	defer dnsManager.Lock.Unlock()
	runtimeGateway := &object.RuntimeGateway{
		Gateway: *gateway,
		Status:  GATEWAY_STATUS_INIT,
	}
	dnsManager.ToBeDoneGatewayMap[gateway.Metadata.Name] = *runtimeGateway
	dnsManager.GatewayMap[gateway.Metadata.Name] = runtimeGateway
	client.AddRuntimeGateway(*runtimeGateway)
}

func deleteGateway(gateway *object.Gateway) {
	serviceList := client.GetServiceByKey(GATEWAY_SERVICE_PREFIX + gateway.Metadata.Name)
	if len(serviceList) == 0 {
		return
	}
	client.DeleteService(serviceList[0])

	replicaSetList := client.GetReplicaSetByKey(GATEWAY_REPLICASET_PREFIX + gateway.Metadata.Name)
	if len(replicaSetList) == 0 {
		return
	}
	client.DeleteReplicaSet(replicaSetList[0])

	dnsManager.Lock.Lock()
	defer dnsManager.Lock.Unlock()
	delete(dnsManager.GatewayMap, gateway.Metadata.Name)
}

func updateGateway(gateway *object.Gateway) {
	deleteGateway(gateway)
	createGateway(gateway)
}
