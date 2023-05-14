package kubeProxy

import (
	"fmt"
	"minik8s/pkg/util/iptables"
)

// DESPERATE

func (rootChain *RootChain) makeRuleCommand() {
	rootChain.RuleCommand = []string{"-j", ROOT_CHAIN, "-s", "0/0", "-d", "0/0", "-p", "all"}
}

func (kubeProxyManager *KubeProxyManager) initRootChain() {
	fmt.Printf("prepare to init root chain\n")
	kubeProxyManager.RootChain.makeRuleCommand()
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
		return
	}
	exist, err2 := ipt.ChainExists(ROOT_TABLE, ROOT_CHAIN)
	if err2 != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	if exist {
		fmt.Printf("root chain exist\n")
		return
	}
	err = ipt.NewChain(ROOT_TABLE, ROOT_CHAIN)
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	err = ipt.Insert(ROOT_TABLE, OUTPUT_CHAIN, 1, kubeProxyManager.RootChain.RuleCommand...)
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	err = ipt.Insert(ROOT_TABLE, PREROUTING_CHAIN, 1, kubeProxyManager.RootChain.RuleCommand...)
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
}

func (kubeProxyManager *KubeProxyManager) deleteRootChain() {
	fmt.Printf("prepare to delete root chain!\n")
	ipt, err := iptables.New()
	if err != nil {
		fmt.Println("delete error")
		fmt.Println(err)
		return
	}
	err = ipt.Delete(ROOT_TABLE, OUTPUT_CHAIN, kubeProxyManager.RootChain.RuleCommand...)
	if err != nil {
		fmt.Println("delete error")
		fmt.Println(err)
	}
	err = ipt.Delete(ROOT_TABLE, PREROUTING_CHAIN, kubeProxyManager.RootChain.RuleCommand...)
	if err != nil {
		fmt.Println("delete error")
		fmt.Println(err)
	}
	for multiServiceKey, multiService := range kubeProxyManager.RootMap {
		for singleServiceKey, singleService := range multiService {
			err := singleService.deleteSingleService()
			if err != nil {
				fmt.Printf("delete error")
			}
			delete(multiService, singleServiceKey)
		}
		delete(kubeProxyManager.RootMap, multiServiceKey)
	}
	err = ipt.DeleteChain(ROOT_TABLE, ROOT_CHAIN)
	if err != nil {
		fmt.Println("delete error")
		fmt.Println(err)
	}
	fmt.Printf("delete root chain done!\n")
}
