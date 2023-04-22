package kubeProxy

import (
	"fmt"
	"minik8s/pkg/service"
)

func Init() {
	fmt.Println("this is proxy")
	service.StartServiceManager()
}
