package kubeProxy

import (
	"bufio"
	"fmt"
	"io"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/services"
	"minik8s/pkg/util/config"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func reloadNginxConfig(name string) {
	res, err := kubelet.GetAllRunningContainers()
	if err != nil {
		fmt.Println("nginx reload fail: " + err.Error())
	}
	var containerIds []string
	for _, val := range res {
		if strings.Contains(val.Names[0], name) {
			containerIds = append(containerIds, val.ID)
		}
	}
	for _, containerId := range containerIds {
		args := fmt.Sprintf("exec %s nginx -s reload", containerId)
		_, err := execCmd("docker", args)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func execCmd(exc string, args string) ([]string, error) {
	cmd := exec.Command(exc, strings.Split(args, " ")...)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(stdout)
	var result []string
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		result = append(result, line)
	}
	err = cmd.Wait()
	return result, err
}

func createDir(path string) {
	args := fmt.Sprintf("-r %s %s", services.NGINX_TEMPLATE_FILEPATH, path)
	fmt.Println("copying nginx template file to " + path)
	res, err := execCmd("cp", args)
	if err != nil {
		fmt.Println("createDir fail")
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func deleteDir(path string) {
	args := fmt.Sprintf("-rf %s", path)
	fmt.Println("removing nginx file in " + path)
	res, err := execCmd("rm", args)
	if err != nil {
		fmt.Println("deleteDir fail")
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func getOriginalDst(clientConn *net.TCPConn) (ipv4 string, port uint16, newTCPConn *net.TCPConn, err error) {

	remoteAddr := clientConn.RemoteAddr()
	if remoteAddr == nil {
		err = fmt.Errorf("clientConn.fd is nil")
		return
	}

	newTCPConn = nil

	clientConnFile, err := clientConn.File()
	if err != nil {
		return
	} else {
		clientConn.Close()
	}

	addr, err := syscall.GetsockoptIPv6Mreq(int(clientConnFile.Fd()), syscall.IPPROTO_IP, 80)
	if err != nil {
		return
	}
	newConn, err := net.FileConn(clientConnFile)
	if err != nil {
		return
	}

	if _, ok := newConn.(*net.TCPConn); ok {
		newTCPConn = newConn.(*net.TCPConn)
		clientConnFile.Close()
	} else {
		fmt.Printf("ERR: newConn is not a *net.TCPConn, instead it is: %T (%v)\n", newConn, newConn)
		return
	}

	ipv4 = strconv.Itoa(int(uint(addr.Multiaddr[4]))) + "." +
		strconv.Itoa(int(uint(addr.Multiaddr[5]))) + "." +
		strconv.Itoa(int(uint(addr.Multiaddr[6]))) + "." +
		strconv.Itoa(int(uint(addr.Multiaddr[7])))
	port = uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])

	return
}

func dial(host string, port int, needTransfer bool) (*net.TCPConn, error) {
	remoteAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nil, err
	}
	remoteAddrAndPort := &net.TCPAddr{IP: remoteAddr.IP, Port: port}
	var localAddrAndPort *net.TCPAddr
	if needTransfer {
		localAddr, err := net.ResolveIPAddr("ip", LOCALHOST_SIDECAR)
		if err != nil {
			return nil, err
		}
		localAddrAndPort = &net.TCPAddr{IP: localAddr.IP, Port: 0}
	} else {
		localAddrAndPort = nil
	}

	conn, err := net.DialTCP("tcp", localAddrAndPort, remoteAddrAndPort)
	return conn, err
}

func copy(dst io.ReadWriteCloser, src io.ReadWriteCloser) {
	if dst == nil || src == nil {
		fmt.Println("[copy] null src/dst")
		return
	}

	defer dst.Close()
	defer src.Close()

	_, err := io.Copy(dst, src)
	if err != nil {
		return
	}
}

func selectorMatch(origin string, selector string, strategy string) bool {
	if strategy == config.VIRTUAL_SERVICE_TYPE_EXACT {
		return origin == selector
	}
	if strategy == config.VIRTUAL_SERVICE_TYPE_PREFIX {
		return strings.HasPrefix(origin, selector)
	}
	if strategy == config.VIRTUAL_SERVICE_TYPE_REGULAR {
		found, err := regexp.MatchString(selector, origin)
		return found && err == nil
	}
	return false
}
