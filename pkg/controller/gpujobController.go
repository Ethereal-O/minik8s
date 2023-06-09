package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

var GpuJobControllerExited = make(chan bool)
var GpuJobControllerToExit = make(chan bool)

func Start_GpuJobController() {
	gpujobChan, stopFunc := messging.Watch("/"+config.GPUJOB_TYPE, true)
	dealJob(gpujobChan)
	fmt.Println("GpuJob Controller start")

	// Wait until Ctrl-C
	<-GpuJobControllerToExit
	stopFunc()
	GpuJobControllerExited <- true
}

func dealJob(gpujobChan chan string) {
	for {
		select {
		case mes := <-gpujobChan:
			if mes == "hello" {
				continue
			}
			//fmt.Println("[this]", mes)
			var tarGpuJob object.GpuJob
			json.Unmarshal([]byte(mes), &tarGpuJob)
			jobname := tarGpuJob.Metadata.Name
			var newPod = &object.Pod{
				Kind: config.POD_TYPE,
				Metadata: object.Metadata{
					Name: object.GpuJobPodFullName(tarGpuJob),
				},
				Spec: object.PodSpec{
					Volumes: []object.Volume{
						{
							Name: "v1",
							Type: "hostPath",
							Path: config.GPU_NODE_DIR_PATH + "/" + jobname,
						},
					},
					Containers: []object.Container{
						{
							Name:  config.GPU_JOB_NAME,
							Image: config.GPU_JOB_IMAGE,
							VolumeMounts: []object.VolumeMount{
								{
									Name:      "v1",
									MountPath: config.GPU_CONTAINER_DIR_PATH,
								},
							},
							Ports: []object.Port{
								{
									ContainerPort: 22,
								},
							},
							Command: []string{config.GPU_JOB_COMMAND, jobname},
							//Command: []string{"/bin/sh", "-c", "touch /tmp/hello.txt;while true;do /bin/echo $(date +%T) >> /tmp/hello.txt;sleep 3;done;"},
						},
					},
				},
			}
			client.AddPod(*newPod)
		}
	}
}
