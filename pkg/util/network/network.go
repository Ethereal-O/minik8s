package network

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex

// GetAvailablePort returns a random available TCP port
func GetAvailablePort() (int, error) {
	// We must use lock here, because many goroutines may call this function
	lock.Lock()
	defer lock.Unlock()
	address, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// IsPortAvailable judges whether given port is available
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("port %s is taken: %s\n", address, err)
		return false
	}

	defer listener.Close()
	return true
}

func GetHostIp() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ipStr := addr.String()
			if strings.HasPrefix(ipStr, "192.168") {
				return strings.Split(ipStr, "/")[0], nil
			}
		}
	}
	return "", errors.New("cannot find IP address which has prefix 192.168")
}

func GetDockerHostIp() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ipStr := addr.String()
			if strings.HasPrefix(ipStr, "172.17") {
				return strings.Split(ipStr, "/")[0], nil
			}
		}
	}
	return "", errors.New("cannot find IP address which has prefix 172.17")
}

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func LookUpIp(hostname string) string {
	addr, _ := net.LookupIP(hostname)
	return fmt.Sprintf("%d.%d.%d.%d/24", addr[0][12], addr[0][13], addr[0][14], addr[0][15])
}
