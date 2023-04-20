package apiServer

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/etcd"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

func unbind(rsName string) {
	pods := etcd.Get_etcd("/"+config.POD_TYPE, true)
	for _, pod := range pods {
		var podObject object.Pod
		json.Unmarshal([]byte(pod), &podObject)
		if podObject.Belong == rsName {
			podObject.Belong = ""
			newPod, err := json.Marshal(podObject)
			if err != nil {
				fmt.Println(err.Error())
			}
			etcd.Set_etcd("/"+config.POD_TYPE+"/"+podObject.Metadata.Name, string(newPod))
		}
	}
}
