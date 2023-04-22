package service

import "minik8s/pkg/object"

func createService(service *object.Service) {
	serviceManager.lock.Lock()

}
