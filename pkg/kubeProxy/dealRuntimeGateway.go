package kubeProxy

import (
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"minik8s/pkg/util/tools"
)

func dealRunningRuntimeGateway(runtimeGateway *object.RuntimeGateway) {
	oldRuntimeGateway, ok := kubeProxyManager.RuntimeGatewayMap[runtimeGateway.Gateway.Metadata.Name]
	if !ok {
		fmt.Printf("creating runtimeGateway %s\n", runtimeGateway.Gateway.Metadata.Name)
		createRuntimeGateway(runtimeGateway)
	} else if tools.MD5(*oldRuntimeGateway) != tools.MD5(*runtimeGateway) {
		fmt.Printf("updating runtimeGateway %s\n", runtimeGateway.Gateway.Metadata.Name)
		updateRuntimeGateway(runtimeGateway)
	} else {
		fmt.Printf("duplicated runtimeGateway %s\n", runtimeGateway.Gateway.Metadata.Name)
	}
}

func dealExitRuntimeGateway(runtimeGateway *object.RuntimeGateway) {
	fmt.Printf("deleting runtimeGateway %s\n", runtimeGateway.Gateway.Metadata.Name)
	deleteRuntimeGateway(runtimeGateway)
}

func createRuntimeGateway(runtimeGateway *object.RuntimeGateway) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	if runtimeGateway.Status == services.GATEWAY_STATUS_INIT {
		createDir(services.GATEWAY_NGINX_PATH_PREFIX + "/" + runtimeGateway.Gateway.Metadata.Name)
	} else if runtimeGateway.Status == services.GATEWAY_STATUS_DEPLOYING {
		runtimeGateway.Status = services.GATEWAY_STATUS_RUNNING
		kubeProxyManager.RuntimeGatewayMap[runtimeGateway.Gateway.Metadata.Name] = runtimeGateway

		client.AddRuntimeGateway(*runtimeGateway)

		updateGatewayNginxConfig(runtimeGateway)
		fmt.Println("write nginx config finished")
		reloadNginxConfig(services.GATEWAY_CONTAINER_PREFIX + runtimeGateway.Gateway.Metadata.Name)
		fmt.Println("reload nginx config finished")
		updateDnsConfig()

	}
}

func deleteRuntimeGateway(runtimeGateway *object.RuntimeGateway) {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	_, ok := kubeProxyManager.RuntimeGatewayMap[runtimeGateway.Gateway.Metadata.Name]
	if !ok {
		return
	}
	delete(kubeProxyManager.RuntimeGatewayMap, runtimeGateway.Gateway.Metadata.Name)
	updateDnsConfig()
	deleteDir(services.GATEWAY_NGINX_PATH_PREFIX + "/" + runtimeGateway.Gateway.Metadata.Name)
}

func updateRuntimeGateway(runtimeGateway *object.RuntimeGateway) {
	if runtimeGateway.Status == services.GATEWAY_STATUS_RUNNING {
		fmt.Printf("not to deal with running gateway")
		return
	}
	deleteRuntimeGateway(runtimeGateway)
	createRuntimeGateway(runtimeGateway)
}
