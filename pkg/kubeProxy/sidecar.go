package kubeProxy

import (
	"fmt"
	"minik8s/pkg/util/network"
	"net"
)

func (kubeProxyManager *KubeProxyManager) initSidecar() error {
	out, err := execCmd("docker", fmt.Sprintf(`run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-iptables -p %v -z %v -u 0 -m REDIRECT -i * -b * -d %v`, OUTBOUND_PORT, INBOUND_PORT, PROMETHEUS_PORT))
	if err != nil {
		fmt.Printf("[InitSidecar] error:%v out:%v\n", err, out)
		fmt.Println("[InitSidecar] Reset iptables...")
		out, err = execCmd("docker", "run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-clean-iptables")
		if err != nil {
			fmt.Printf("[InitSidecar] Reset fail, error:%v out:%v\n", err, out)
			return err
		}
		out, err = execCmd("docker", fmt.Sprintf(`run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-iptables -p %v -z %v -u 0 -m REDIRECT -i * -b * -d %v`, OUTBOUND_PORT, INBOUND_PORT, PROMETHEUS_PORT))
		return err
	}

	go kubeProxyManager.listenSocket(LOCALHOST, OUTBOUND_PORT)
	go kubeProxyManager.listenSocket(LOCALHOST, INBOUND_PORT)
	localhost_docker, err := network.GetDockerHostIp()
	if err != nil {
		go kubeProxyManager.listenSocket(localhost_docker, OUTBOUND_PORT)
		go kubeProxyManager.listenSocket(localhost_docker, INBOUND_PORT)
	} else {
		go kubeProxyManager.listenSocket(LOCALHOST_DOCKER, OUTBOUND_PORT)
		go kubeProxyManager.listenSocket(LOCALHOST_DOCKER, INBOUND_PORT)
	}

	return err
}

func (kubeProxyManager *KubeProxyManager) deleteSidecar() error {
	out, err := execCmd("docker", "run --rm --name=istio-init --network=host --cap-add=NET_ADMIN istio/proxyv2:1.16.0 istio-clean-iptables")
	if err != nil {
		fmt.Printf("[DeleteSidecar] error:%v out:%v\n", err, out)
	}
	return err
}

func (kubeProxyManager *KubeProxyManager) listenSocket(ip string, port string) {
	lnaddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("[ListenSocket] No port available")
		return
	}
	server, err := net.ListenTCP("tcp", lnaddr)
	if err != nil {
		fmt.Println("[listenSocket] init server fail")
		return
	}
	fmt.Println("[listenSocket] listening to " + ip + ":" + port)
	if port == INBOUND_PORT && ip == LOCALHOST {
		for {
			conn, _ := server.AcceptTCP()
			go dealConn(conn, true)
		}
	} else {
		for {
			conn, _ := server.AcceptTCP()
			go dealConn(conn, false)
		}
	}
}

func dealConn(clientConn *net.TCPConn, needTransfer bool) {
	if clientConn == nil {
		return
	}

	ipv4, port, clientConn, err := getOriginalDst(clientConn)

	if err != nil {
		return
	}

	targetIP := getTargetIp(ipv4)

	directConn, err := dial(targetIP, int(port), needTransfer)
	if err != nil {
		fmt.Printf("Could not connect, giving up: %v", err)
		return
	}

	go copy(directConn, clientConn)
	go copy(clientConn, directConn)
}
