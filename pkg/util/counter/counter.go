package counter

import (
	"fmt"
	"math/big"
	"minik8s/pkg/etcd"
	"net"
	"strconv"
	"sync"
)

func deserializeIP(ip int) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func serializeIP(ip string) string {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return strconv.Itoa(int(ret.Int64()))
}

var uuidCounter = Counter{
	url:       UUIDCounterURL,
	initCount: UUIDCounterInitCount,
}

var monitorCounter = Counter{
	url:       MonitorCounterURL,
	initCount: MonitorCounterInitCount,
}

var rrPolicyCounter = Counter{
	url:       RRPolicyCounterURL,
	initCount: RRPolicyCounterInitCount,
}

var nodePortCounter = Counter{
	url:       NodePortCounterURL,
	initCount: NodePortCounterInitPort,
}

var podIPCounter = Counter{
	url:       PodIPCounterURL,
	initCount: serializeIP(PodIPCounterInitIP),
}

var serviceIPCounter = Counter{
	url:       ServiceIPCounterURL,
	initCount: serializeIP(ServiceIPCounterInitIP),
}

var nodeIPCounter = Counter{
	url:       NodeIPCounterURL,
	initCount: serializeIP(NodeIPCounterInitIP),
}

func GetUuid() string {
	return strconv.Itoa(uuidCounter.fetchAndAdd())
}

func GetMonitorCrt() string {
	return strconv.Itoa(monitorCounter.fetchAndAdd())
}

func GetRRPolicy() int {
	return rrPolicyCounter.fetchAndAdd()
}

func NewNodePort() string {
	return strconv.Itoa(nodePortCounter.fetchAndAdd())
}

func NewPodIP() string {
	return deserializeIP(podIPCounter.fetchAndAdd())
}

func NewServiceIP() string {
	return deserializeIP(serviceIPCounter.fetchAndAdd())
}

func NewNodeIP() string {
	return deserializeIP(nodeIPCounter.fetchAndAdd())
}

type Counter struct {
	url       string
	initCount string
	lock      sync.Mutex
}

func (counter *Counter) fetchAndAdd() int {
	counter.lock.Lock()
	defer counter.lock.Unlock()

	ret := etcd.Get_etcd(counter.url, false)

	if len(ret) != 1 {
		etcd.Set_etcd(counter.url, counter.initCount)
		ret[0] = counter.initCount
	}

	num, _ := strconv.Atoi(ret[0])
	etcd.Set_etcd(counter.url, strconv.Itoa(num+1))
	return num
}
