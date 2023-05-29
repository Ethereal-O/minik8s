package kubeProxy

import (
	"minik8s/pkg/client"
	"minik8s/pkg/services"
	"minik8s/pkg/util/config"
)

func (kubeProxyManager *KubeProxyManager) initKubeProxyManager() {
	kubeProxyManager.Lock.Lock()
	defer kubeProxyManager.Lock.Unlock()
	kubeProxyManager.initGatewayMap()
	kubeProxyManager.initServiceMap()
	updateDnsConfig()
	if config.SERVICE_POLICY == config.SERVICE_POLICY_NGINX {
		applyNodePortService()
	}
	if config.SERVICE_POLICY == config.SERVICE_POLICY_IPTABLES {
		applyAllMultiService()
	}
}

func (kubeProxyManager *KubeProxyManager) initServiceMap() {
	allRuntimeServices := client.GetAllRuntimeServices()
	for _, runtimeService := range allRuntimeServices {
		if runtimeService.Status == services.SERVICE_STATUS_RUNNING {
			runtimeServiceRef := runtimeService
			kubeProxyManager.RuntimeServiceMap[runtimeService.Service.Metadata.Name] = &runtimeServiceRef
		}
	}
}

func (kubeProxyManager *KubeProxyManager) initGatewayMap() {
	allRuntimeGateways := client.GetAllRuntimeGateways()
	for _, runtimeGateway := range allRuntimeGateways {
		if runtimeGateway.Status == services.GATEWAY_STATUS_RUNNING || runtimeGateway.Status == services.GATEWAY_STATUS_DEPLOYING {
			runtimeGatewayRef := runtimeGateway
			runtimeGatewayRef.Status = services.GATEWAY_STATUS_RUNNING
			kubeProxyManager.RuntimeGatewayMap[runtimeGateway.Gateway.Metadata.Name] = &runtimeGatewayRef
		}
	}
}
