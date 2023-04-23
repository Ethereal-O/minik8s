package kubeProxy

import (
	"fmt"
	"minik8s/pkg/service"
)

func Start_proxy() {
	fmt.Println("this is proxy")
	service.StartServiceManager()
}
