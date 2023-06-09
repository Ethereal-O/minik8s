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

// --------------------------- RuntimeService ---------------------------

func GetAllRuntimeServices() []object.RuntimeService {
	runtimeServiceList := Get_object(config.EMPTY_FLAG, config.RUNTIMESERVICE_TYPE)
	var resList []object.RuntimeService
	for _, runtimeService := range runtimeServiceList {
		var runtimeServiceObject object.RuntimeService
		json.Unmarshal([]byte(runtimeService), &runtimeServiceObject)
		resList = append(resList, runtimeServiceObject)
	}
	return resList
}

func GetRuntimeServiceByKey(key string) []object.RuntimeService {
	runtimeServiceList := Get_object(key, config.RUNTIMESERVICE_TYPE)
	var resList []object.RuntimeService
	for _, runtimeService := range runtimeServiceList {
		var runtimeServiceObject object.RuntimeService
		json.Unmarshal([]byte(runtimeService), &runtimeServiceObject)
		resList = append(resList, runtimeServiceObject)
	}
	return resList
}

func AddRuntimeService(runtimeService object.RuntimeService) string {
	serviceStatusValue, err := json.Marshal(runtimeService)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(runtimeService.Service.Metadata.Name, string(serviceStatusValue), config.RUNTIMESERVICE_TYPE)
}

func DeleteRuntimeService(runtimeService object.RuntimeService) string {
	return Delete_object(runtimeService.Service.Metadata.Name, config.RUNTIMESERVICE_TYPE)
}

// --------------------------- Virtual Service ---------------------------

func GetAllVirtualServices() []object.VirtualService {
	virtualServiceList := Get_object(config.EMPTY_FLAG, config.VIRTUALSERVICE_TYPE)
	var resList []object.VirtualService
	for _, virtualService := range virtualServiceList {
		var virtualServiceObject object.VirtualService
		json.Unmarshal([]byte(virtualService), &virtualServiceObject)
		resList = append(resList, virtualServiceObject)
	}
	return resList
}

func GetVirtualServiceByKey(key string) []object.VirtualService {
	virtualServiceList := Get_object(key, config.VIRTUALSERVICE_TYPE)
	var resList []object.VirtualService
	for _, virtualService := range virtualServiceList {
		var virtualServiceObject object.VirtualService
		json.Unmarshal([]byte(virtualService), &virtualServiceObject)
		resList = append(resList, virtualServiceObject)
	}
	return resList
}

func AddVirtualService(virtualService object.VirtualService) string {
	virtualServiceValue, err := json.Marshal(virtualService)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(virtualService.Metadata.Name, string(virtualServiceValue), config.VIRTUALSERVICE_TYPE)
}

func DeleteVirtualService(virtualService object.VirtualService) string {
	return Delete_object(virtualService.Metadata.Name, config.VIRTUALSERVICE_TYPE)
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
	return Put_object(gateway.Metadata.Name, string(gatewayValue), config.GATEWAY_TYPE)
}

func DeleteGateway(gateway object.Gateway) string {
	return Delete_object(gateway.Metadata.Name, config.GATEWAY_TYPE)
}

// --------------------------- RuntimeGateway ---------------------------

func GetAllRuntimeGateways() []object.RuntimeGateway {
	runtimegatewayList := Get_object(config.EMPTY_FLAG, config.RUNTIMEGATEWAY_TYPE)
	var resList []object.RuntimeGateway
	for _, runtimeGateway := range runtimegatewayList {
		var runtimeGatewayObject object.RuntimeGateway
		json.Unmarshal([]byte(runtimeGateway), &runtimeGatewayObject)
		resList = append(resList, runtimeGatewayObject)
	}
	return resList
}

func AddRuntimeGateway(runtimeGateway object.RuntimeGateway) string {
	gatewayStatusValue, err := json.Marshal(runtimeGateway)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(runtimeGateway.Gateway.Metadata.Name, string(gatewayStatusValue), config.RUNTIMEGATEWAY_TYPE)
}

func DeleteRuntimeGateway(runtimeGateway object.RuntimeGateway) string {
	return Delete_object(runtimeGateway.Gateway.Metadata.Name, config.RUNTIMEGATEWAY_TYPE)
}

// --------------------------- ReplicaSet ---------------------------

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

// --------------------------- DaemonSet ---------------------------

func AddDaemonSet(daemonSet object.DaemonSet) string {
	daemonSetValue, err := json.Marshal(daemonSet)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(daemonSet.Metadata.Name, string(daemonSetValue), config.DAEMONSET_TYPE)
}

func GetDaemonSetByKey(key string) []object.DaemonSet {
	daemonSetList := Get_object(key, config.DAEMONSET_TYPE)
	var resList []object.DaemonSet
	for _, daemonSet := range daemonSetList {
		var daemonSetObject object.DaemonSet
		json.Unmarshal([]byte(daemonSet), &daemonSetObject)
		resList = append(resList, daemonSetObject)
	}
	return resList
}

func DeleteDaemonSet(daemonSet object.DaemonSet) string {
	return Delete_object(daemonSet.Metadata.Name, config.DAEMONSET_TYPE)
}

// --------------------------- AutoScaler ---------------------------

func AddAutoScaler(autoScaler object.AutoScaler) string {
	autoScalerValue, err := json.Marshal(autoScaler)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(autoScaler.Metadata.Name, string(autoScalerValue), config.AUTOSCALER_TYPE)
}

func GetAllAutoScalers() []object.AutoScaler {
	hpaList := Get_object(config.EMPTY_FLAG, config.AUTOSCALER_TYPE)
	var resList []object.AutoScaler
	for _, hpa := range hpaList {
		var hpaObject object.AutoScaler
		json.Unmarshal([]byte(hpa), &hpaObject)
		resList = append(resList, hpaObject)
	}
	return resList
}

func DeleteAutoScaler(autoScaler object.AutoScaler) string {
	return Delete_object(autoScaler.Metadata.Name, config.AUTOSCALER_TYPE)
}

// --------------------------- Node ---------------------------

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

func GetActiveNodes() []object.Node {
	nodeList := Get_object(config.EMPTY_FLAG, config.NODE_TYPE)
	var resList []object.Node
	for _, node := range nodeList {
		var nodeObject object.Node
		json.Unmarshal([]byte(node), &nodeObject)
		if nodeObject.Runtime.Status != config.EXIT_STATUS {
			resList = append(resList, nodeObject)
		}
	}
	return resList
}

func GetNode(key string) object.Node {
	nodeList := Get_object(key, config.NODE_TYPE)
	var nodeObject object.Node
	json.Unmarshal([]byte(nodeList[0]), &nodeObject)
	return nodeObject
}

func AddNode(node object.Node) string {
	nodeValue, err := json.Marshal(node)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(node.Metadata.Name, string(nodeValue), config.NODE_TYPE)
}

// --------------------------- GpuJob ---------------------------

func AddGpuJob(gpuJob object.GpuJob) string {
	gpuJobValue, err := json.Marshal(gpuJob)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(gpuJob.Metadata.Name, string(gpuJobValue), config.GPUJOB_TYPE)
}

func GetAllGpuJob() []object.GpuJob {
	gpuJobList := Get_object(config.EMPTY_FLAG, config.GPUJOB_TYPE)
	var resList []object.GpuJob
	for _, gpuJob := range gpuJobList {
		var gpuJobObject object.GpuJob
		json.Unmarshal([]byte(gpuJob), &gpuJobObject)
		resList = append(resList, gpuJobObject)
	}
	return resList
}

func DeleteGpuJob(gpuJob object.GpuJob) string {
	return Delete_object(gpuJob.Metadata.Name, config.GPUJOB_TYPE)
}

// --------------------------- ServerlessFunctions ---------------------------

func AddServerlessFunctions(serverlessFunctions object.ServerlessFunctions) string {
	serverlessFunctionsValue, err := json.Marshal(serverlessFunctions)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(serverlessFunctions.Metadata.Name, string(serverlessFunctionsValue), config.SERVERLESSFUNCTIONS_TYPE)
}

func GetAllServerlessFunctions() []object.ServerlessFunctions {
	serverlessFunctionsList := Get_object(config.EMPTY_FLAG, config.SERVERLESSFUNCTIONS_TYPE)
	var resList []object.ServerlessFunctions
	for _, serverlessFunctions := range serverlessFunctionsList {
		var serverlessFunctionsObject object.ServerlessFunctions
		json.Unmarshal([]byte(serverlessFunctions), &serverlessFunctionsObject)
		resList = append(resList, serverlessFunctionsObject)
	}
	return resList
}

func DeleteServerlessFunctions(serverlessFunctions object.ServerlessFunctions) string {
	return Delete_object(serverlessFunctions.Metadata.Name, config.SERVERLESSFUNCTIONS_TYPE)
}

func GetAllFunctions() []object.Function {
	var functionList []object.Function
	serverlessFunctionsList := GetAllServerlessFunctions()
	for _, serverlessFunctions := range serverlessFunctionsList {
		for _, function := range serverlessFunctions.Spec.Items {
			function.Runtime = serverlessFunctions.Runtime
			function.FaasName = serverlessFunctions.Metadata.Name
			functionList = append(functionList, function)
		}
	}
	return functionList
}

func GetActiveFunctions() []object.Function {
	var functionList []object.Function
	serverlessFunctionsList := GetAllServerlessFunctions()
	for _, serverlessFunctions := range serverlessFunctionsList {
		if serverlessFunctions.Runtime.Status != config.RUNNING_STATUS {
			continue
		}
		for _, function := range serverlessFunctions.Spec.Items {
			function.Runtime = serverlessFunctions.Runtime
			function.FaasName = serverlessFunctions.Metadata.Name
			functionList = append(functionList, function)
		}
	}
	return functionList
}

func GetFunction(funcName string) *object.Function {
	functionList := GetAllFunctions()
	for _, function := range functionList {
		if function.FuncName == funcName {
			return &function
		}
	}
	return nil
}
