package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"time"
)

var FaasControllerExited = make(chan bool)
var FaasControllerToExit = make(chan bool)

var FaasExited = make(map[string]chan bool)
var FaasToExit = make(map[string]chan bool)

func Start_serverlessFunctionsController() {
	funcChan, stopFunc := messging.Watch("/"+config.SERVERLESSFUNCTIONS_TYPE, true)
	dealFunc(funcChan)
	fmt.Println("ServerlessFunctions Controller start")

	// Wait until Ctrl-C
	<-FaasControllerToExit
	stopFunc()
	FaasControllerExited <- true
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
			if tarFuncSet.Runtime.Status == "" {
				createRs(tarFuncSet)
				time.Sleep(10 * time.Second)
				createService(tarFuncSet)
				FaasToExit[tarFuncSet.Metadata.Name] = make(chan bool)
				FaasExited[tarFuncSet.Metadata.Name] = make(chan bool)
				go FaasCycle(tarFuncSet)
			} else if tarFuncSet.Runtime.Status == config.EXIT_STATUS {
				FaasToExit[tarFuncSet.Metadata.Name] <- true
				<-FaasExited[tarFuncSet.Metadata.Name]
				delete(FaasToExit, tarFuncSet.Metadata.Name)
				delete(FaasExited, tarFuncSet.Metadata.Name)
				tarRs := client.GetReplicaSetByKey(object.FaasRsFullName(tarFuncSet))[0]
				client.DeleteReplicaSet(tarRs)
				tarService := client.GetServiceByKey(object.FaasServiceFullName(tarFuncSet))[0]
				client.DeleteService(tarService)
			}
		}
	}
}

func createRs(tarFuncSet object.ServerlessFunctions) {
	name := tarFuncSet.Metadata.Name
	var newRs = &object.ReplicaSet{
		Kind: config.REPLICASET_TYPE,
		Metadata: object.Metadata{
			Name: object.FaasRsFullName(tarFuncSet),
			Labels: map[string]string{
				"faas": object.FaasRsFullName(tarFuncSet),
			},
		},
		Spec: object.RsSpec{
			Replicas: 0,
			Template: object.Template{
				Metadata: object.Metadata{
					Name: object.FaasPodFullName(tarFuncSet),
					Labels: map[string]string{
						"faas": object.FaasPodFullName(tarFuncSet),
					},
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
			},
		},
	}
	client.AddReplicaSet(*newRs)
}
func createService(tarFuncSet object.ServerlessFunctions) {
	var newService = &object.Service{
		Kind: config.SERVICE_TYPE,
		Metadata: object.Metadata{
			Name: object.FaasServiceFullName(tarFuncSet),
		},
		Spec: object.ServiceSpec{
			Type: "ClusterIP",
			Ports: []object.ServicePort{
				{
					Port:       "8081",
					TargetPort: "8081",
					Protocol:   "tcp",
				},
			},
			Selector: map[string]string{
				"faas": object.FaasPodFullName(tarFuncSet),
			},
		},
	}
	client.AddService(*newService)
}

func FaasCycle(tarFuncSet object.ServerlessFunctions) {
	serviceName := object.FaasServiceFullName(tarFuncSet)
	serviceChan, stopFunc := messging.Watch("/"+config.RUNTIMESERVICE_TYPE+"/"+serviceName, false)

	for {
		select {
		case <-FaasToExit[tarFuncSet.Metadata.Name]:
			FaasExited[tarFuncSet.Metadata.Name] <- true
			stopFunc()
			return
		case mes := <-serviceChan:
			var tarService object.RuntimeService
			json.Unmarshal([]byte(mes), &tarService)
			if len(tarService.Pods) > 0 {
				tarFuncSet.Runtime.Status = config.RUNNING_STATUS
				tarFuncSet.Runtime.FunctionIp = tarService.Service.Runtime.ClusterIp
				client.AddServerlessFunctions(tarFuncSet)
			} else {
				tarFuncSet.Runtime.Status = config.CREATED_STATUS
				client.AddServerlessFunctions(tarFuncSet)
			}
		}
	}
}
