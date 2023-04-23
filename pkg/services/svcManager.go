package services

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"sync"
)

var serviceManager *ServiceManager
var Exited = make(chan bool)
var ToExit = make(chan bool)

func createServiceManager() *ServiceManager {
	serviceManager := &ServiceManager{}
	serviceManager.ServiceMap = make(map[string]ServiceStatus)
	serviceManager.StopChannel = make(chan struct{})
	serviceManager.DnsMap = make(map[string]DnsStatus)
	var lock sync.Mutex
	serviceManager.Lock = lock
	return serviceManager
}

func StartServiceManager() {
	serviceManager = createServiceManager()
	serviceChan, serviceStop := messging.Watch("/"+config.SERVICE_TYPE, true)
	go dealService(serviceChan)

	// Wait until Ctrl-C
	<-ToExit
	serviceStop()
	Exited <- true
}

//func startTicker() {
//    //ticker := time.NewTicker(5 * time.Second)
//    //go func() {
//    //	for {
//    //		select {
//    //		case <-ticker.C:
//    //			serviceManager.checkAndBoot()
//    //		case <-serviceManager.stopChannel:
//    //			ticker.Stop()
//    //			return
//    //		}
//    //	}
//    //}()
//}

//func (serviceManager *ServiceManager) dnsLoop() {
//	for {
//		time.Sleep(2 * time.Second)
//		serviceManager.lock.Lock()
//		var removes []string
//		for k, v := range serviceManager.DnsMap {
//			clien
//			resp, err := serviceManager.client.GetRuntimeService(netconfig.GateWayServicePrefix + k)
//			if err != nil {
//				fmt.Println("[checkDnsAndTrans] getRuntimeService fail" + err.Error())
//				continue
//			}
//			if resp == nil {
//				continue
//			}
//			if resp.Status.Phase == object.Running {
//				v.Status.Phase = object.ServiceCreated
//				v.Spec.GateWayIp = resp.Spec.ClusterIp
//				err = serviceManager.client.UpdateDnsAndTrans(v)
//				if err != nil {
//					fmt.Println("[checkDnsAndService]updateDns fail" + err.Error())
//					continue
//				}
//				removes = append(removes, k)
//			}
//		}
//		for _, val := range removes {
//			delete(serviceManager.name2DnsMap, val)
//		}
//		serviceManager.lock.Unlock()
//	}
//}
//
//// 没隔一段时间查看一下有无节点注册， 如果有注册的调用boot
//func (manager *Manager) checkAndBoot() {
//	for {
//		time.Sleep(5 * time.Second)
//		res, err := manager.ls.List(config.NODE_PREFIX)
//		if err != nil {
//			fmt.Println("[ServiceManager] checkAndBoot error" + err.Error())
//			continue
//		}
//		if len(res) == 0 {
//			continue
//		} else {
//			manager.boot()
//			break
//		}
//	}
//}
//func (manager *Manager) boot() {
//	//查看一下是否已经存在coreDns service, 存在的话不再生成
//	res, err := manager.client.GetRuntimeService("dnsService")
//	if res != nil {
//		return
//	}
//	//生成coreDns service
//	err = manager.client.AddConfigRs(GetCoreDnsRsModule())
//	if err != nil {
//		fmt.Println("[ServiceManager] boot fail" + err.Error())
//		return
//	}
//	time.Sleep(1 * time.Second)
//	err = manager.client.UpdateService(GetCoreDnsServiceModule())
//	if err != nil {
//		fmt.Println("[ServiceManager] boot fail" + err.Error())
//		return
//	}
//}

func dealService(serviceChan chan string) {
	for {
		select {
		case mes := <-serviceChan:
			var tarService object.Service
			err := json.Unmarshal([]byte(mes), &tarService)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			if tarService.Runtime.Status == config.EXIT_STATUS {
				dealExitService(&tarService)
			} else {
				dealRunningService(&tarService)
			}
		}
	}
}
