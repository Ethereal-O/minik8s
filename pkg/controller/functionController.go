package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

var ServerlessFunctionsControllerExited = make(chan bool)
var ServerlessFunctionsControllerToExit = make(chan bool)

func Start_serverlessFunctionsController() {
	funcChan, stopFunc := messging.Watch("/"+config.SERVERLESSFUNCTIONS_TYPE, true)
	dealFunc(funcChan)
	fmt.Println("GpuJob Controller start")

	// Wait until Ctrl-C
	<-ServerlessFunctionsControllerToExit
	stopFunc()
	ServerlessFunctionsControllerExited <- true
}

func dealFunc(funcChan chan string) {
	for {
		select {
		case mes := <-funcChan:
			if mes == "hello" {
				continue
			}
			//fmt.Println("[this]", mes)
			var tarFuncSet object.ServerlessFunctions
			json.Unmarshal([]byte(mes), &tarFuncSet)
			name := tarFuncSet.Metadata.Name
			var newPod = &object.Pod{
				Kind: config.POD_TYPE,
				Metadata: object.Metadata{
					Name: "faas_pod_" + name,
				},
				Spec: object.PodSpec{
					Volumes: []object.Volume{
						{
							Name: "v1",
							Type: "hostPath",
							Path: config.FUNC_NODE_DIR_PATH + "/" + name,
						},
					},
					Containers: []object.Container{
						{
							Name:  config.FUNC_NAME,
							Image: config.FUNC_IMAGE,
							VolumeMounts: []object.VolumeMount{
								{
									Name:      "v1",
									MountPath: config.FUNC_CONTAINER_DIR_PATH,
								},
							},
							Ports: []object.Port{
								{
									ContainerPort: 22,
								},
							},
							Command: []string{config.FUNC_COMMAND},
						},
					},
				},
			}
			client.AddPod(*newPod)
		}
	}
}
