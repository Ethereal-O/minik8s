package kubeProxy

import (
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/util/iptables"
)

func createSingleService(runtimeService *object.RuntimeService, port object.ServicePort, podsInfo []PodInfo) *SingleService {
	singleService := &SingleService{
		Table:  PARENT_TABLE,
		Parent: PARENT_CHAIN,
		// we can specify a single service by its name and port
		Name:        SINGLE_SERVICE + runtimeService.Service.Metadata.Name + "-" + port.Port,
		ClusterPort: port.Port,
		ClusterIp:   runtimeService.Service.Spec.ClusterIp,
		Protocol:    port.Protocol,
	}
	singleService.makeRuleCommand()
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("make new ipTable error")
		fmt.Println(err)
		return nil
	}
	err = ipt.NewChain(singleService.Table, singleService.Name)
	i := 0
	singleService.SinglePodMap = make(map[string]*SinglePod)
	for _, podInfo := range podsInfo {
		singlePod := createSinglePod(singleService, podInfo, len(podsInfo)-i)
		err = singlePod.initSinglePod()
		if err != nil {
			fmt.Println("make singlePod error")
			fmt.Println(err)
			return nil
		}
		singleService.SinglePodMap[podInfo.PodName] = singlePod
		i++
	}
	return singleService
}

func (singleService *SingleService) makeRuleCommand() {
	singleService.RuleCommand = []string{"-s", "0/0", "-d", singleService.ClusterIp, "-p", singleService.Protocol, "--dport", singleService.ClusterPort, "-j", singleService.Name}
}

func (singleService *SingleService) initSingleService() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Append(singleService.Table, singleService.Parent, singleService.RuleCommand...)
	if err != nil {
		return err
	}
	return nil
}

func (singleService *SingleService) deleteSingleService() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Delete(singleService.Table, singleService.Parent, singleService.RuleCommand...)
	if err != nil {
		return err
	}
	for _, singlePod := range singleService.SinglePodMap {
		err = singlePod.deleteSinglePod()
		if err != nil {
			return err
		}
	}
	err = ipt.DeleteChain(singleService.Table, singleService.Name)
	return err
}
