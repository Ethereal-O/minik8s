package services

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"sync"
)

func dealRunningService(service *object.Service) {
	_, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		createService(service)
	} else {
		updateService(service)
	}
}

func dealExitService(service *object.Service) {
	deleteService(service)
}

func createService(service *object.Service) {
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	serviceStatus := object.ServiceStatus{
		Service: *service,
		Lock:    sync.Mutex{},
		Pods:    []object.Pod{},
	}
	InitServiceStatus(&serviceStatus)
	serviceManager.ServiceMap[service.Metadata.Name] = serviceStatus
}

func deleteService(service *object.Service) {
	serviceManager.Lock.Lock()
	defer serviceManager.Lock.Unlock()
	serviceStatus, ok := serviceManager.ServiceMap[service.Metadata.Name]
	if !ok {
		return
	}
	serviceStatus.Timer.Stop()
	ret := client.DeleteServiceStatus(serviceStatus)
	fmt.Println(ret)
	delete(serviceManager.ServiceMap, service.Metadata.Name)
}

func updateService(service *object.Service) {
	deleteService(service)
	createService(service)
}
