package kubeProxy

import "minik8s/pkg/util/iptables"

func createSingleNet(singlePod *SinglePod, podInfo PodInfo) *SingleNet {
	singleNet := &SingleNet{
		Table:    singlePod.Table,
		Parent:   singlePod.Name,
		Name:     SINGLE_NET + podInfo.PodName + podInfo.PodPort,
		PodIp:    podInfo.PodIP,
		Port:     podInfo.PodPort,
		Protocol: singlePod.Protocol,
	}
	singleNet.makeRuleCommand()
	return singleNet
}

func (singleNet *SingleNet) makeRuleCommand() {
	singleNet.RuleCommand = []string{"-s", "0/0", "-d", "0/0", "-p", singleNet.Protocol, "-j", "DNAT", "--to-destination", singleNet.PodIp + ":" + singleNet.Port}
}

func (singleNet *SingleNet) initSingleNet() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Append(singleNet.Table, singleNet.Parent, singleNet.RuleCommand...)
	if err != nil {
		return err
	}
	return nil
}

func (singleNet *SingleNet) deleteSingleNet() error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}
	err = ipt.Delete(singleNet.Table, singleNet.Parent, singleNet.RuleCommand...)
	if err != nil {
		return err
	}
	return nil
}
