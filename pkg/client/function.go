package client

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
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

func GetAllServiceStatuses() []services.ServiceStatus {
	serviceStatusList := Get_object(config.EMPTY_FLAG, config.SERVICESTATUS_TYPE)
	var resList []services.ServiceStatus
	for _, serviceStatus := range serviceStatusList {
		var serviceStatusObject services.ServiceStatus
		json.Unmarshal([]byte(serviceStatus), &serviceStatusObject)
		resList = append(resList, serviceStatusObject)
	}
	return resList
}

func AddServiceStatus(serviceStatus services.ServiceStatus) string {
	serviceStatusValue, err := json.Marshal(serviceStatus)
	if err != nil {
		fmt.Println(err.Error())
	}
	return Put_object(serviceStatus.Service.Metadata.Name, string(serviceStatusValue), config.SERVICESTATUS_TYPE)
}

func DeleteServiceStatus(serviceStatus services.ServiceStatus) string {
	return Delete_object(serviceStatus.Service.Metadata.Name, config.SERVICESTATUS_TYPE)
}
