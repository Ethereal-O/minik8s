package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/iptables"
)

// DESPERATE

func createSingleService(runtimeService *object.RuntimeService, port object.ServicePort, podsInfo []PodInfo) *SingleService {
	singleService := &SingleService{
		Table:  ROOT_TABLE,
		Parent: ROOT_CHAIN,
		// we can specify a single service by its name and port
		Name:        SINGLE_SERVICE + "-" + runtimeService.Service.Metadata.Name + "-" + port.Port,
		ClusterPort: port.Port,
		ClusterIp:   runtimeService.Service.Runtime.ClusterIp,
		IsNodePort:  runtimeService.Service.Spec.Type == config.SERVICE_TYPE_NODEPORT,
		NodePort:    port.NodePort,
		Protocol:    port.Protocol,
	}
	singleService.makeRuleCommand()
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("make new ipTable error")
		fmt.Println(err.Error())
		return singleService
	}
	err = ipt.NewChain(singleService.Table, singleService.Name)
	i := 0
	singleService.SinglePodMap = make(map[string]*SinglePod)
	for _, podInfo := range podsInfo {
		singlePod := createSinglePod(singleService, podInfo, len(podsInfo)-i)
		err = singlePod.initSinglePod()
		if err != nil {
			fmt.Println("make singlePod error")
			fmt.Println(err.Error())
		}
		singleService.SinglePodMap[podInfo.PodName] = singlePod
		i++
	}
	return singleService
}

func (singleService *SingleService) makeRuleCommand() {
	singleService.RuleCommandClusterIp = []string{"-s", "0/0", "-d", singleService.ClusterIp, "-p", singleService.Protocol, "--dport", singleService.ClusterPort, "-j", singleService.Name}
	singleService.RuleCommandNodePort = []string{"-s", "0/0", "-p", singleService.Protocol, "--dport", singleService.NodePort, "-j", singleService.Name}
}

func (singleService *SingleService) initSingleService() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Append(singleService.Table, singleService.Parent, singleService.RuleCommandClusterIp...)
	if singleService.IsNodePort {
		err = ipt.Append(singleService.Table, singleService.Parent, singleService.RuleCommandNodePort...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (singleService *SingleService) deleteSingleService() error {
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ipt.Delete(singleService.Table, singleService.Parent, singleService.RuleCommandClusterIp...)
	if singleService.IsNodePort {
		err = ipt.Delete(singleService.Table, singleService.Parent, singleService.RuleCommandNodePort...)
	}
	for _, singlePod := range singleService.SinglePodMap {
		err = singlePod.deleteSinglePod()
	}
	err = ipt.DeleteChain(singleService.Table, singleService.Name)
	return err
}
