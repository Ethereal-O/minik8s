package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var dnsManager *DnsManager
var dnsManagerExited = make(chan bool)
var dnsManagerToExit = make(chan os.Signal)

func createDnsManager() *DnsManager {
	dnsManager := &DnsManager{}
	dnsManager.GatewayMap = make(map[string]object.RuntimeGateway)
	var lock sync.Mutex
	dnsManager.Lock = lock
	return dnsManager
}

func StartDnsManager() {
	dnsManager = createDnsManager()
	dnsManager.initDnsManager()
	signal.Notify(dnsManagerToExit, syscall.SIGINT, syscall.SIGTERM)
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
