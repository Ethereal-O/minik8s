package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"sync"
)

var dnsManager *DnsManager

func createDnsManager() *DnsManager {
	dnsManager := &DnsManager{}
	dnsManager.ToBeDoneGatewayMap = make(map[string]object.RuntimeGateway)
	dnsManager.GatewayMap = make(map[string]*object.RuntimeGateway)
	var lock sync.Mutex
	dnsManager.Lock = lock
	return dnsManager
}

func StartDnsManager() {
	dnsManager = createDnsManager()
	dnsManager.initDnsManager()
	gatewayChan, dnsStop := messging.Watch("/"+config.GATEWAY_TYPE, true)
	go dealGateway(gatewayChan)

	// Wait until Ctrl-C
	<-ToExit
	// please note: dnsManager is not have duty to delete dns service and replica set, otherwise it will cause double delete
	// serviceManager and replicaSetManager will manage the lifecycle of dns service and replica set
	dnsStop()
	Exited <- true
}

func dealGateway(gatewayChan chan string) {
	for {
		select {
		case mes := <-gatewayChan:
			if mes == "hello" {
				continue
			}
			var tarGateway object.Gateway
			err := json.Unmarshal([]byte(mes), &tarGateway)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarGateway.Runtime.Status == config.EXIT_STATUS {
				dealExitGateway(&tarGateway)
			} else if tarGateway.Runtime.Status == config.RUNNING_STATUS {
				dealRunningGateway(&tarGateway)
			} else {
				fmt.Println("Gateway status error!")
			}
		}
	}
}
