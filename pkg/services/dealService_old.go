package services

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/tools"
	"sync"
)

func dealRunningService_old(service *object.Service) {
	oldRuntimeService, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		fmt.Printf("creating service %s\n", service.Metadata.Name)
		createService_old(service)
	} else if tools.MD5(oldRuntimeService.Service) != tools.MD5(*service) {
		fmt.Printf("updating service %s\n", service.Metadata.Name)
		updateService_old(service)
	} else {
		fmt.Printf("duplicated service %s\n", service.Metadata.Name)
	}
}

func dealExitService_old(service *object.Service) {
	fmt.Printf("deleting service %s\n", service.Metadata.Name)
	deleteService_old(service)
}

func createService_old(service *object.Service) {
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	runtimeService := &object.RuntimeService{
		Service: *service,
		Status:  SERVICE_STATUS_RUNNING,
		Lock:    sync.Mutex{},
		Pods:    []object.Pod{},
	}
	InitRuntimeService_old(runtimeService)
	serviceManager.ServiceMap[service.Metadata.Name] = runtimeService
}

func deleteService_old(service *object.Service) {
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	runtimeService, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		return
	}
	runtimeService.Timer.Stop()
	ret := client.DeleteRuntimeService(*runtimeService)
	fmt.Println(ret)
	delete(serviceManager.ServiceMap, service.Metadata.Name)
}

func updateService_old(service *object.Service) {
	deleteService(service)
	createService_old(service)
}
