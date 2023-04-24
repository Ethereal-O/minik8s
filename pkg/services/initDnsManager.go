package services

import (
	"minik8s/pkg/client"
	"minik8s/pkg/util/config"
	"time"
)

func (dnsManager *DnsManager) initDnsManager() {
	dnsManager.Timer = *time.NewTicker(CHECK_DNS_TIME_INTERVAL)
	dnsManager.InitDnsTemplate()
	go dnsManager.checkDnsLoop()
	go dnsManager.checkGatewayLoop()
}

func (dnsManager *DnsManager) checkDnsLoop() {
	defer dnsManager.Timer.Stop()
	for {
		select {
		case <-dnsManager.Timer.C:
			dnsManager.checkDns()
		}
	}
}

func (dnsManager *DnsManager) checkGatewayLoop() {
	defer dnsManager.Timer.Stop()
	for {
		select {
		case <-dnsManager.Timer.C:
			dnsManager.checkGateway()
		}
	}
}

func (dnsManager *DnsManager) checkDns() {
	nodeRes := client.GetAllNodes()
	if len(nodeRes) == 0 {
		return
	}
	dnsRes := client.GetRuntimeServiceByKey(config.DNS_TYPE)
	// if existed, don't need to create
	if dnsRes != nil {
		return
	}
	// create dns replica set
	client.AddReplicaSet(GetDnsReplicaSet())
	time.Sleep(CREATE_RS_AND_SERVICE_TIME_INTERVAL)
	// create dns service
	client.AddService(GetDnsService())
}

func (dnsManager *DnsManager) checkGateway() {
	dnsManager.Lock.Lock()
	defer dnsManager.Lock.Unlock()
	dnsManager.transferGatewayToKubeProxy()
}

func (dnsManager *DnsManager) transferGatewayToKubeProxy() {
	// wait until service online and transfer to kube-proxy
	var removes []string
	for gatewayName, gateWayStatus := range dnsManager.GatewayMap {
		resList := client.GetRuntimeServiceByKey(GATEWAY_SERVICE_PREFIX + gatewayName)
		if len(resList) == 0 {
			continue
		}
		gateWayStatus.Status = GATEWAY_STATUS_DEPLOYING
		gateWayStatus.ClusterIp = resList[0].Service.Spec.ClusterIp
		client.AddRuntimeGateway(gateWayStatus)
		removes = append(removes, gatewayName)
	}
	for _, gateway := range removes {
		delete(dnsManager.GatewayMap, gateway)
	}
}
