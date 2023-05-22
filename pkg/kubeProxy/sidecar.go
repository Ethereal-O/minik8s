package kubeProxy

import (
	"fmt"
	"net"
)

func (kubeProxyManager *KubeProxyManager) initSidecar() error {
	out, err := execCmd("docker", `run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-iptables -p 15001 -z 15006 -u 1337 -m REDIRECT -i * -b * -d 15020,4161`)
	if err != nil {
		fmt.Printf("[InitSidecar] error:%v out:%v\n", err, out)
		fmt.Println("[InitSidecar] Reset iptables...")
		out, err = execCmd("docker", "run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-clean-iptables")
		if err != nil {
			fmt.Printf("[InitSidecar] Reset fail, error:%v out:%v\n", err, out)
			return err
		}
		out, err = execCmd("docker", `run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-iptables -p 15001 -z 15006 -u 1337 -m REDIRECT -i * -b * -d 15020,4161`)
		return err
	}
	go kubeProxyManager.listenSocket()
	return err
}

func (kubeProxyManager *KubeProxyManager) deleteSidecar() error {
	out, err := execCmd("docker", "run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-clean-iptables")
	if err != nil {
		fmt.Printf("[DeleteSidecar] error:%v out:%v\n", err, out)
	}
	return err
}

func (kubeProxyManager *KubeProxyManager) listenSocket() {
	lnaddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:15001")
	if err != nil {
		fmt.Println("[ListenSocket] No port available")
		return
	}
	server, err := net.ListenTCP("tcp", lnaddr)
	if err != nil {
		fmt.Println("[listenSocket] init server fail")
		return
	}
	fmt.Println("[listenSocket] listening to " + "127.0.0.1:15001")
	for {
		conn, _ := server.AcceptTCP()
		ipv4, port, _, _ := getOriginalDst(conn)
		fmt.Println(ipv4)
		fmt.Println(port)
	}
}
