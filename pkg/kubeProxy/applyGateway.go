package kubeProxy

import (
	"bufio"
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"minik8s/pkg/util/config"
	"os"
)

func updateDnsConfig() {
	f, err := os.OpenFile(services.HOST_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("dns config write fail")
		return
	}
	w := bufio.NewWriter(f)
	var hosts_str string
	for _, gateway := range kubeProxyManager.RuntimeGatewayMap {
		if gateway.Status != services.GATEWAY_STATUS_RUNNING {
			continue
		}
		hosts_str += fmt.Sprintf("%s %s\n", gateway.ClusterIp, gateway.Gateway.Spec.Host)
	}
	hosts_str += fmt.Sprintf("%s %s", config.PIP3_SOURCE_IMAGE_IP, config.PIP3_SOURCE_IMAGE_HOSTNAME)
	_, err = fmt.Fprintln(w, hosts_str)
	err = w.Flush()
	if err != nil {
		fmt.Println("dns config write fail: " + err.Error())
	}

	// write it to host
	f_bak, err := os.OpenFile(services.HOST_HOSTS_BAK_PATH, os.O_RDONLY, 0644)
	defer f_bak.Close()
	// read all data to string
	var str string
	data := make([]byte, 1024)
	for {
		n, err := f_bak.Read(data)
		if err != nil {
			break
		}
		str += string(data[:n])
	}
	str += "\n"
	str += hosts_str
	f_host, err := os.OpenFile(services.HOST_HOSTS_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f_host.Close()
	_, err = fmt.Fprintln(f_host, str)
	return
}

func updateGatewayNginxConfig(runtimeGateway *object.RuntimeGateway) {
	var content []string
	content = append(content, makeGatewayConfig(runtimeGateway)...)
	f, err := os.OpenFile(services.GATEWAY_NGINX_PATH_PREFIX+"/"+runtimeGateway.Gateway.Metadata.Name+"/"+services.NGINX_CONFIG_FILE, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("nginx config write fail: " + err.Error())
		return
	}
	w := bufio.NewWriter(f)
	for _, v := range content {
		fmt.Fprintln(w, v)
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("nginx config write fail: " + err.Error())
		return
	}
	return
}

func makeGatewayConfig(runtimeGateway *object.RuntimeGateway) []string {
	var result []string
	result = append(result, "error_log stderr;")
	result = append(result, "events { worker_connections  1024; }")
	result = append(result, "http {", "    access_log /dev/stdout combined;")
	result = append(result, "    server {", "        listen 80 ;")
	result = append(result, fmt.Sprintf("        server_name %s;", runtimeGateway.Gateway.Spec.Host))
	for _, path := range runtimeGateway.Gateway.Spec.Paths {
		result = append(result, fmt.Sprintf("        location %s {", path.Name))
		result = append(result, fmt.Sprintf("            proxy_pass http://%s:%s/;", path.IP, path.Port))
		result = append(result, "        }")
	}
	result = append(result, "       }")
	result = append(result, "}")
	return result
}
