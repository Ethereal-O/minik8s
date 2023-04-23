package client

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func Get_object(key string, tp string) []string {
	url := config.APISERVER_URL
	for _, conftp := range config.TP {
		if tp == conftp {
			url += "/" + conftp + "/" + key
			return get(url)
		}
	}
	return nil
}

func Put_object(key string, value string, tp string) string {
	url := config.APISERVER_URL
	for _, conftp := range config.TP {
		if tp == conftp {
			url += "/" + conftp + "/" + key
			return put(url, value)
		}
	}
	return "not found such type in Put_object!"
}

func Delete_object(key string, tp string) string {
	url := config.APISERVER_URL
	for _, conftp := range config.TP {
		if tp == conftp {
			url += "/" + conftp + "/" + key
			return delete(url)
		}
	}
	return "not found such type in Delete_object!"
}

func Post(key string, prix bool, crt string) string {
	return postFormData(key, prix, crt)
}

// --------------------------- Pod ---------------------------

func GetAllPods() []object.Pod {
	podList := Get_object(config.EMPTY_FLAG, config.POD_TYPE)
	var resList []object.Pod
	for _, pod := range podList {
		var podObject object.Pod
		json.Unmarshal([]byte(pod), &podObject)
		resList = append(resList, podObject)
	}
	return resList
}

func GetActivePods() []object.Pod {
	podList := Get_object(config.EMPTY_FLAG, config.POD_TYPE)
	var resList []object.Pod
	for _, pod := range podList {
		var podObject object.Pod
		json.Unmarshal([]byte(pod), &podObject)
		if podObject.Runtime.Status != config.EXIT_STATUS {
			resList = append(resList, podObject)
		}
	}
	return resList
}

func AddPod(pod object.Pod) string {
	podValue, err := json.Marshal(pod)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(pod.Metadata.Name, string(podValue), config.POD_TYPE)
}

func DeletePod(pod object.Pod) string {
	return Delete_object(pod.Metadata.Name, config.POD_TYPE)
}

// --------------------------- Service ---------------------------

func GetAllServices() []object.Service {
	serviceList := Get_object(config.EMPTY_FLAG, config.SERVICE_TYPE)
	var resList []object.Service
	for _, service := range serviceList {
		var serviceObject object.Service
		json.Unmarshal([]byte(service), &serviceObject)
		resList = append(resList, serviceObject)
	}
	return resList
}

func GetServiceByKey(key string) []object.Service {
	serviceList := Get_object(key, config.SERVICE_TYPE)
	var resList []object.Service
	for _, service := range serviceList {
		var serviceObject object.Service
		json.Unmarshal([]byte(service), &serviceObject)
		resList = append(resList, serviceObject)
	}
	return resList
}

func AddService(service object.Service) string {
	serviceValue, err := json.Marshal(service)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(service.Metadata.Name, string(serviceValue), config.SERVICE_TYPE)
}

func DeleteService(service object.Service) string {
	return Delete_object(service.Metadata.Name, config.SERVICE_TYPE)
}

// --------------------------- ServiceStatus ---------------------------

func GetAllServiceStatuses() []object.ServiceStatus {
	serviceStatusList := Get_object(config.EMPTY_FLAG, config.SERVICESTATUS_TYPE)
	var resList []object.ServiceStatus
	for _, serviceStatus := range serviceStatusList {
		var serviceStatusObject object.ServiceStatus
		json.Unmarshal([]byte(serviceStatus), &serviceStatusObject)
		resList = append(resList, serviceStatusObject)
	}
	return resList
}

func GetServiceStatusByKey(key string) []object.ServiceStatus {
	serviceStatusList := Get_object(key, config.SERVICESTATUS_TYPE)
	var resList []object.ServiceStatus
	for _, serviceStatus := range serviceStatusList {
		var serviceStatusObject object.ServiceStatus
		json.Unmarshal([]byte(serviceStatus), &serviceStatusObject)
		resList = append(resList, serviceStatusObject)
	}
	return resList
}

func AddServiceStatus(serviceStatus object.ServiceStatus) string {
	serviceStatusValue, err := json.Marshal(serviceStatus)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(serviceStatus.Service.Metadata.Name, string(serviceStatusValue), config.SERVICESTATUS_TYPE)
}

func DeleteServiceStatus(serviceStatus object.ServiceStatus) string {
	return Delete_object(serviceStatus.Service.Metadata.Name, config.SERVICESTATUS_TYPE)
}

// --------------------------- Gateway ---------------------------

func GetAllGateways() []object.Gateway {
	gatewayList := Get_object(config.EMPTY_FLAG, config.GATEWAY_TYPE)
	var resList []object.Gateway
	for _, gateway := range gatewayList {
		var gatewayObject object.Gateway
		json.Unmarshal([]byte(gateway), &gatewayObject)
		resList = append(resList, gatewayObject)
	}
	return resList
}

func AddGateway(gateway object.Gateway) string {
	gatewayValue, err := json.Marshal(gateway)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(gateway.MetaData.Name, string(gatewayValue), config.GATEWAY_TYPE)
}

func DeleteGateway(gateway object.Gateway) string {
	return Delete_object(gateway.MetaData.Name, config.GATEWAY_TYPE)
}

// --------------------------- GatewayStatus ---------------------------

func GetAllGatewayStatuses() []object.GatewayStatus {
	gatewayStatusList := Get_object(config.EMPTY_FLAG, config.GATEWAYSTATUS_TYPE)
	var resList []object.GatewayStatus
	for _, gatewayStatus := range gatewayStatusList {
		var gatewayStatusObject object.GatewayStatus
		json.Unmarshal([]byte(gatewayStatus), &gatewayStatusObject)
		resList = append(resList, gatewayStatusObject)
	}
	return resList
}

func AddGatewayStatus(gatewayStatus object.GatewayStatus) string {
	gatewayStatusValue, err := json.Marshal(gatewayStatus)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(gatewayStatus.Gateway.MetaData.Name, string(gatewayStatusValue), config.GATEWAYSTATUS_TYPE)
}

func DeleteGatewayStatus(gatewayStatus object.GatewayStatus) string {
	return Delete_object(gatewayStatus.Gateway.MetaData.Name, config.GATEWAYSTATUS_TYPE)
}

// --------------------------- ReplicaSet ---------------------------

// TODO: add ReplicaSet

func AddReplicaSet(replicaSet object.ReplicaSet) string {
	replicaSetValue, err := json.Marshal(replicaSet)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(replicaSet.Metadata.Name, string(replicaSetValue), config.REPLICASET_TYPE)
}

func GetReplicaSetByKey(key string) []object.ReplicaSet {
	replicaSetList := Get_object(key, config.REPLICASET_TYPE)
	var resList []object.ReplicaSet
	for _, replicaSet := range replicaSetList {
		var replicaSetObject object.ReplicaSet
		json.Unmarshal([]byte(replicaSet), &replicaSetObject)
		resList = append(resList, replicaSetObject)
	}
	return resList
}

func DeleteReplicaSet(replicaSet object.ReplicaSet) string {
	return Delete_object(replicaSet.Metadata.Name, config.REPLICASET_TYPE)
}

// --------------------------- Nodes ---------------------------

// TODO: add Node

func GetAllNodes() []object.Node {
	nodeList := Get_object(config.EMPTY_FLAG, config.NODE_TYPE)
	var resList []object.Node
	for _, node := range nodeList {
		var nodeObject object.Node
		json.Unmarshal([]byte(node), &nodeObject)
		resList = append(resList, nodeObject)
	}
	return resList
}
