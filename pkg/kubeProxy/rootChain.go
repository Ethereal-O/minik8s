package kubeProxy

import (
	"fmt"
	"minik8s/pkg/util/iptables"
)

func (kubeProxyManager *KubeProxyManager) initRootChain() {
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
		return
	}
	err = ipt.NewChain(ROOT_TABLE, ROOT_CHAIN)
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	err = ipt.Insert(ROOT_TABLE, OUTPUT_CHAIN, 1, "-j", ROOT_CHAIN, "-s", "0/0", "-d", "0/0", "-p", "all")
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
	err = ipt.Insert(ROOT_TABLE, PREROUTING_CHAIN, 1, "-j", ROOT_CHAIN, "-s", "0/0", "-d", "0/0", "-p", "all")
	if err != nil {
		fmt.Println("boot error")
		fmt.Println(err)
	}
}

func (kubeProxyManager *KubeProxyManager) deleteRootChain() {
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
}
