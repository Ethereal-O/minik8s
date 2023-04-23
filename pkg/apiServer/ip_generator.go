package apiServer

import (
	"fmt"
	"math/big"
	"minik8s/pkg/etcd"
	"minik8s/pkg/util/network"
	"net"
	"strconv"
)

var podIPGenerator = ipGenerator{
	url:    network.PodIPGeneratorURL,
	initIp: network.PodIPGeneratorInitIP,
}

func NewPodIP() string {
	return podIPGenerator.AddAndFetch()
}

func deserialize(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func serialize(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

type ipGenerator struct {
	url    string
	initIp string
}

func (ig *ipGenerator) AddAndFetch() string {
	ret := etcd.Get_etcd(ig.url, false)
	num, _ := strconv.Atoi(ret[0])
	etcd.Set_etcd(ig.url, strconv.Itoa(num+1))
	return deserialize(int64(num + 1))
}

func (ig *ipGenerator) Init() {
	ret := etcd.Get_etcd(ig.url, false)
	if len(ret) != 1 {
		etcd.Set_etcd(ig.url, strconv.Itoa(int(serialize(ig.initIp))))
	}
}
