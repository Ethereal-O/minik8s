package kubeProxy

import (
	"bufio"
	"fmt"
	"io"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"os"
	"os/exec"
	"strings"
)

func updateDnsConfig() {
	f, err := os.OpenFile(services.HOST_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("dns config write fail")
		return
	}
	w := bufio.NewWriter(f)
	for _, gateway := range kubeProxyManager.RuntimeGatewayMap {
		lineStr := fmt.Sprintf("%s %s", gateway.ClusterIp, gateway.Gateway.Spec.Host)
		_, err := fmt.Fprintln(w, lineStr)
		if err != nil {
			fmt.Println("dns config write fail: " + err.Error())
			return
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("dns config write fail: " + err.Error())
	}
	return
}

func updateNginxConfig(runtimeGateway *object.RuntimeGateway) {
	var content []string
	content = append(content, makeConfig(runtimeGateway)...)
	f, err := os.OpenFile(services.NGINX_PATH_PREFIX+"/"+runtimeGateway.Gateway.Metadata.Name+"/"+services.NGINX_CONFIG_FILE, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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

func reloadNginxConfig(runtimeGateway *object.RuntimeGateway) {
	res, err := kubelet.GetAllRunningContainers()
	if err != nil {
		fmt.Println("nginx reload fail: " + err.Error())
	}
	var containerIds []string
	for _, val := range res {
		if strings.Contains(val.Names[0], services.GATEWAY_CONTAINER_PREFIX+runtimeGateway.Gateway.Metadata.Name) {
			containerIds = append(containerIds, val.ID)
		}
	}
	for _, containerId := range containerIds {
		args := fmt.Sprintf("exec %s nginx -s reload", containerId)
		res, err := execCmd("docker", args)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
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

func makeConfig(runtimeGateway *object.RuntimeGateway) []string {
	var result []string
	result = append(result, "error_log stderr;")
	result = append(result, "events { worker_connections  1024; }")
	result = append(result, "http {", "    access_log /dev/stdout combined;")
	result = append(result, "    server {", "        listen 80 ;")
	result = append(result, fmt.Sprintf("        server_name %s;", runtimeGateway.Gateway.Spec.Host))
	for _, path := range runtimeGateway.Gateway.Spec.Paths {
		result = append(result, fmt.Sprintf("        location %s {", path.Name))
		result = append(result, fmt.Sprintf("            proxy_pass http://%s:%s;", path.IP, path.Port))
		result = append(result, "        }")
	}
	result = append(result, "       }")
	result = append(result, "}")
	return result
}

func createDir(runtimeGateway *object.RuntimeGateway) {
	args := fmt.Sprintf("-r %s %s", services.NGINX_TEMPLATE_FILEPATH, services.NGINX_PATH_PREFIX+"/"+runtimeGateway.Gateway.Metadata.Name)
	res, err := execCmd("cp", args)
	if err != nil {
		fmt.Println("createDir fail")
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func deleteDir(runtimeGateway *object.RuntimeGateway) {
	args := fmt.Sprintf("-rf %s", services.NGINX_PATH_PREFIX+"/"+runtimeGateway.Gateway.Metadata.Name)
	res, err := execCmd("rm", args)
	if err != nil {
		fmt.Println("deleteDir fail")
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
