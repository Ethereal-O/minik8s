package kubeProxy

import (
	"fmt"
	"minik8s/pkg/util/iptables"
)

func createSinglePod(singleService *SingleService, podInfo PodInfo, id int) *SinglePod {
	singlePod := &SinglePod{
		Table:    singleService.Table,
		Parent:   singleService.Name,
		Name:     SINGLE_POD + podInfo.PodName + podInfo.PodPort,
		Protocol: singleService.Protocol,
		PodName:  podInfo.PodName,
		Id:       id,
	}
	singlePod.makeRuleCommand()
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("make new ipTable error")
		fmt.Println(err)
		return nil
	}
	err = ipt.NewChain(singlePod.Table, singlePod.Name)
	if err != nil {
		fmt.Println("make new chain error")
		fmt.Println(err)
		return nil
	}
	singlePod.Net = createSingleNet(singlePod, podInfo)
	err = singlePod.Net.initSingleNet()
	if err != nil {
		fmt.Println("init singleNet error")
		fmt.Println(err)
	}
	return singlePod
}

func (singlePod *SinglePod) makeRuleCommand() {
	singlePod.RuleCommand = []string{"-p", singlePod.Protocol, "-m", "statistic", "--mode", "nth", "--every", fmt.Sprintf("%d", singlePod.Id), "--packet", "0", "-j", singlePod.Name}
}

func (singlePod *SinglePod) initSinglePod() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Append(singlePod.Table, singlePod.Parent, singlePod.RuleCommand...)
	if err != nil {
		return err
	}
	return nil
}

func (singlePod *SinglePod) deleteSinglePod() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Delete(singlePod.Table, singlePod.Parent, singlePod.RuleCommand...)
	if err != nil {
		return err
	}
	err = singlePod.Net.deleteSingleNet()
	if err != nil {
		return err
	}
	err = ipt.DeleteChain(singlePod.Table, singlePod.Name)
	return err
}
