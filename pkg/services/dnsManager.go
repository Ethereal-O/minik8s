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
var dnsManagerExited = make(chan bool)
var dnsManagerToExit = make(chan bool)

func createDnsManager() *DnsManager {
	dnsManager := &DnsManager{}
	dnsManager.GatewayMap = make(map[string]object.GatewayStatus)
	var lock sync.Mutex
	dnsManager.Lock = lock
	return dnsManager
}

func StartDnsManager() {
	dnsManager = createDnsManager()
	dnsManager.InitDnsManager()
	gatewayChan, dnsStop := messging.Watch("/"+config.GATEWAY_TYPE, true)
	go dealGateway(gatewayChan)

	// Wait until Ctrl-C
	<-dnsManagerToExit
	dnsStop()
	dnsManagerExited <- true
}

func dealGateway(gatewayChan chan string) {
	for {
		select {
		case mes := <-gatewayChan:
			var tarGateway object.Gateway
			err := json.Unmarshal([]byte(mes), &tarGateway)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarGateway.Runtime.Status == config.EXIT_STATUS {
				dealExitGateway(&tarGateway)
			} else {
				dealRunningGateway(&tarGateway)
			}
		}
	}
}
