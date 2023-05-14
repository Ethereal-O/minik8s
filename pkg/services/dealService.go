package services

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
	"sync"
)

func dealRunningService(service *object.Service) {
	oldRuntimeService, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		fmt.Printf("creating service %s\n", service.Metadata.Name)
		createService(service)
	} else if tools.MD5(oldRuntimeService.Service) != tools.MD5(*service) {
		fmt.Printf("updating service %s\n", service.Metadata.Name)
		updateService(service)
	} else {
		fmt.Printf("duplicated service %s\n", service.Metadata.Name)
	}
}

func dealExitService(service *object.Service) {
	fmt.Printf("deleting service %s\n", service.Metadata.Name)
	deleteService(service)
}

func createService(service *object.Service) {
	client.AddReplicaSet(GetServiceReplicaSet(service.Metadata.Name))
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	runtimeService := &object.RuntimeService{
		Service: *service,
		Status:  SERVICE_STATUS_INIT,
		Lock:    sync.Mutex{},
		Pods:    []object.Pod{},
	}
	client.AddRuntimeService(*runtimeService)
	runtimeService.Status = SERVICE_STATUS_RUNNING
	InitRuntimeService(runtimeService)
	serviceManager.ServiceMap[service.Metadata.Name] = runtimeService
}

func deleteService(service *object.Service) {
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	runtimeService, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		return
	}
	runtimeService.Timer.Stop()
	ret := client.DeleteRuntimeService(*runtimeService)
	fmt.Println(ret)

	replicaSetList := client.GetReplicaSetByKey(SERVICE_REPLICASET_PREFIX + service.Metadata.Name)
	if len(replicaSetList) == 0 {
		return
	}
	client.DeleteReplicaSet(replicaSetList[0])

	delete(serviceManager.ServiceMap, service.Metadata.Name)
}

func updateService(service *object.Service) {
	deleteService(service)
	createService(service)
}
