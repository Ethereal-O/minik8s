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

var ServerlessFunctionsControllerExited = make(chan bool)
var ServerlessFunctionsControllerToExit = make(chan bool)

func Start_serverlessFunctionsController() {
	funcChan, stopFunc := messging.Watch("/"+config.SERVERLESSFUNCTIONS_TYPE, true)
	dealFunc(funcChan)
	fmt.Println("ServerlessFunctions Controller start")

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
			if tarFuncSet.Runtime.Status == "" || tarFuncSet.Runtime.Status == config.CREATED_STATUS {
				//createPod(tarFuncSet)

				createRs(tarFuncSet)
				time.Sleep(10 * time.Second)
				createService(tarFuncSet)
			}
		}
	}
}

func createRs(tarFuncSet object.ServerlessFunctions) {
	name := tarFuncSet.Metadata.Name
	var newRs = &object.ReplicaSet{
		Kind: config.REPLICASET_TYPE,
		Metadata: object.Metadata{
			Name: object.ServerlessFunctionsRsFullName(tarFuncSet),
			Labels: map[string]string{
				"faas": object.ServerlessFunctionsRsFullName(tarFuncSet),
			},
		},
		Spec: object.RsSpec{
			Replicas: 3,
			Template: object.Template{
				Metadata: object.Metadata{
					Name: object.ServerlessFunctionsPodFullName(tarFuncSet),
					Labels: map[string]string{
						"faas": object.ServerlessFunctionsPodFullName(tarFuncSet),
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
func createHpa(tarFuncSet object.ServerlessFunctions) {
	var newHpa = &object.AutoScaler{
		Kind: config.AUTOSCALER_TYPE,
		Metadata: object.Metadata{
			Name: object.ServerlessFunctionsHpaFullName(tarFuncSet),
		},
		Spec: object.HpaSpec{
			MinReplicas: 3,
			MaxReplicas: 3,
			Interval:    30,
			ScaleTargetRef: object.TargetRef{
				Kind: config.REPLICASET_TYPE,
				Name: object.ServerlessFunctionsRsFullName(tarFuncSet),
			},
		},
	}
	client.AddAutoScaler(*newHpa)
}
func createService(tarFuncSet object.ServerlessFunctions) {
	var newService = &object.Service{
		Kind: config.SERVICE_TYPE,
		Metadata: object.Metadata{
			Name: object.ServerlessFunctionsServiceFullName(tarFuncSet),
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
				"faas": object.ServerlessFunctionsPodFullName(tarFuncSet),
			},
		},
	}
	client.AddService(*newService)
	go dealFaasService(tarFuncSet, object.ServerlessFunctionsServiceFullName(tarFuncSet))
}

func dealFaasService(tarFuncSet object.ServerlessFunctions, servicename string) {
	serviceChan, stopFunc := messging.Watch("/"+config.RUNTIMESERVICE_TYPE+"/"+servicename, false)

	for {
		select {
		case mes := <-serviceChan:
			var tarService object.RuntimeService
			json.Unmarshal([]byte(mes), &tarService)
			if len(tarService.Pods) > 0 {
				tarFuncSet.Runtime.Status = config.RUNNING_STATUS
				tarFuncSet.Runtime.FunctionIp = tarService.Service.Runtime.ClusterIp
				client.AddServerlessFunctions(tarFuncSet)
			}
		}
	}
	stopFunc()
}

// ---------------------------------------------------------
func createPod(tarFuncSet object.ServerlessFunctions) {
	name := tarFuncSet.Metadata.Name
	var newPod = &object.Pod{
		Kind: config.POD_TYPE,
		Metadata: object.Metadata{
			Name: object.ServerlessFunctionsPodFullName(tarFuncSet),
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
	go dealFaasPod(tarFuncSet, object.ServerlessFunctionsPodFullName(tarFuncSet))
}

func dealFaasPod(tarFuncSet object.ServerlessFunctions, podname string) {
	podChan, stopFunc := messging.Watch("/"+config.POD_TYPE+"/"+podname, false)

	for {
		select {
		case mes := <-podChan:
			var tarPod object.Pod
			json.Unmarshal([]byte(mes), &tarPod)
			if tarFuncSet.Runtime.Status != tarPod.Runtime.Status || tarFuncSet.Runtime.FunctionIp != tarPod.Runtime.ClusterIp {
				tarFuncSet.Runtime.Status = tarPod.Runtime.Status
				tarFuncSet.Runtime.FunctionIp = tarPod.Runtime.ClusterIp
				client.AddServerlessFunctions(tarFuncSet)
			}

		}
	}
	stopFunc()
}
