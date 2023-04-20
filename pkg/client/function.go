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

// ------------------------------------------------------------------------------------------

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

func GetRunningPods() []object.Pod {
	podList := Get_object(config.EMPTY_FLAG, config.POD_TYPE)
	var resList []object.Pod
	for _, pod := range podList {
		var podObject object.Pod
		json.Unmarshal([]byte(pod), &podObject)
		if podObject.Metadata.Status == config.RUNNING_STATUS {
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
