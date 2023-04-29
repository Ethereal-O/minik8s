package apiServer

import (
	"minik8s/pkg/etcd"
	"minik8s/pkg/util/network"
	"strconv"
)

var nodePortGenerator = portGenerator{
	url:      network.NodePortGeneratorURL,
	initPort: network.NodePortGeneratorInitPort,
}

func NewNodePort() string {
	return nodePortGenerator.FetchAndAdd()
}

type portGenerator struct {
	url      string
	initPort string
}

func (pg *portGenerator) FetchAndAdd() string {
	ret := etcd.Get_etcd(pg.url, false)
	num, _ := strconv.Atoi(ret[0])
	etcd.Set_etcd(pg.url, strconv.Itoa(num+1))
	return ret[0]
}

func (pg *portGenerator) Init() {
	ret := etcd.Get_etcd(pg.url, false)
	if len(ret) != 1 {
		etcd.Set_etcd(pg.url, pg.initPort)
	}
}
