package kubeProxy

import (
	"bufio"
	"fmt"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/network"
	"minik8s/pkg/util/tools"
	"minik8s/pkg/util/weave"
	"os"
	"strings"
)

func applyService(runtimeService *object.RuntimeService) {
	applyClusterIpService(runtimeService)
	if runtimeService.Service.Spec.Type != config.SERVICE_TYPE_NODEPORT {
		return
	}
	applyNodePortService()
}

func applyClusterIpService(runtimeService *object.RuntimeService) {
	applyWeaveAttach(runtimeService)
	updateServiceNginxConfig(runtimeService)
	fmt.Println("write nginx config finished")
	reloadNginxConfig(services.SERVICE_CONTAINER_PREFIX + runtimeService.Service.Metadata.Name)
	fmt.Println("reload nginx config finished")
}

func applyNodePortService() {
	updateNodePortServiceNginxConfig()
	fmt.Println("write nginx config finished")
	reloadNginxConfig(services.FORWARD_DAEMONSET_PREFIX)
	fmt.Println("reload nginx config finished")
}

func updateServiceNginxConfig(runtimeService *object.RuntimeService) {
	var content []string
	content = append(content, makeServiceConfig(runtimeService)...)
	f, err := os.OpenFile(services.SERVICE_NGINX_PATH_PREFIX+"/"+runtimeService.Service.Metadata.Name+"/"+services.NGINX_CONFIG_FILE, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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

func makeServiceConfig(runtimeService *object.RuntimeService) []string {
	var result []string
	result = append(result, "error_log stderr;")
	result = append(result, "events { worker_connections  1024; }")

	// http block
	result = append(result, "http {", "    access_log /dev/stdout combined;")
	for _, port := range runtimeService.Service.Spec.Ports {
		if port.Protocol == "UDP" {
			continue
		}
		result = append(result, fmt.Sprintf("    upstream %s {", runtimeService.Service.Metadata.Name+"-"+port.Port))
		for _, pod := range runtimeService.Pods {
			result = append(result, fmt.Sprintf("        server %s:%s;", pod.Runtime.ClusterIp, port.TargetPort))
		}
		result = append(result, "    }")
	}
	for _, port := range runtimeService.Service.Spec.Ports {
		if port.Protocol == "UDP" {
			continue
		}
		result = append(result, fmt.Sprintf("    server {        listen %s ;", port.Port))
		result = append(result, fmt.Sprintf("        server_name localhost;"))
		result = append(result, fmt.Sprintf("        location / {"))
		result = append(result, fmt.Sprintf("            proxy_pass http://%s/;", runtimeService.Service.Metadata.Name+"-"+port.Port))
		result = append(result, "        }")
		result = append(result, "       }")
	}
	result = append(result, "}")

	// udp block
	result = append(result, "stream {")
	for _, port := range runtimeService.Service.Spec.Ports {
		if port.Protocol == "TCP" {
			continue
		}
		result = append(result, fmt.Sprintf("    upstream %s {", runtimeService.Service.Metadata.Name+"-"+port.Port))
		for _, pod := range runtimeService.Pods {
			result = append(result, fmt.Sprintf("        server %s:%s;", pod.Runtime.ClusterIp, port.TargetPort))
		}
		result = append(result, "    }")
	}
	for _, port := range runtimeService.Service.Spec.Ports {
		if port.Protocol == "TCP" {
			continue
		}
		result = append(result, fmt.Sprintf("    server {        listen %s udp ;", port.Port))
		result = append(result, fmt.Sprintf("            proxy_pass %s;", runtimeService.Service.Metadata.Name+"-"+port.Port))
		result = append(result, "       }")
	}
	result = append(result, "}")

	return result
}

func applyWeaveAttach(runtimeService *object.RuntimeService) {
	// get all pods and selector
	allPods := client.GetAllPods()
	// first check if service-pod is running
	runningPods, _ := tools.Filter(allPods, func(pod object.Pod) bool {
		if pod.Runtime.Status == config.RUNNING_STATUS && strings.Contains(pod.Metadata.Name, runtimeService.Service.Metadata.Name) {
			return true
		} else {
			return false
		}
	})

	if len(runningPods) == 0 {
		fmt.Printf("this should not happen, service %s has no running pod\n", runtimeService.Service.Metadata.Name)
		return
	}

	// attach clusterIp to pod
	err := weave.Attach(runningPods[0].Runtime.Containers[0], runtimeService.Service.Runtime.ClusterIp+network.Mask)
	if err != nil {
		fmt.Println(err)
		//return
	}
}

func updateNodePortServiceNginxConfig() {
	var content []string
	content = append(content, makeNodePortServiceConfig()...)
	f, err := os.OpenFile(services.FORWARD_NGINX_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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

func makeNodePortServiceConfig() []string {
	var result []string
	result = append(result, "error_log stderr;")
	result = append(result, "events { worker_connections  1024; }")

	// http block
	result = append(result, "http {", "    access_log /dev/stdout combined;")
	for _, runtimeService := range kubeProxyManager.RuntimeServiceMap {
		if runtimeService.Service.Spec.Type != config.SERVICE_TYPE_NODEPORT {
			continue
		}
		for _, port := range runtimeService.Service.Spec.Ports {
			if port.Protocol == "UDP" {
				continue
			}
			result = append(result, fmt.Sprintf("    server {        listen %s ;", port.NodePort))
			result = append(result, fmt.Sprintf("        server_name localhost;"))
			result = append(result, fmt.Sprintf("        location / {"))
			result = append(result, fmt.Sprintf("            proxy_pass http://%s:%s/;", runtimeService.Service.Runtime.ClusterIp, port.Port))
			result = append(result, "        }")
			result = append(result, "       }")
		}
	}
	result = append(result, "}")

	// udp block
	result = append(result, "stream {")
	for _, runtimeService := range kubeProxyManager.RuntimeServiceMap {
		if runtimeService.Service.Spec.Type != config.SERVICE_TYPE_NODEPORT {
			continue
		}
		for _, port := range runtimeService.Service.Spec.Ports {
			if port.Protocol == "TCP" {
				continue
			}
			result = append(result, fmt.Sprintf("    server {        listen %s udp ;", port.NodePort))
			result = append(result, fmt.Sprintf("            proxy_pass %s:%s;", runtimeService.Service.Runtime.ClusterIp, port.Port))
			result = append(result, "       }")
		}
	}
	result = append(result, "}")

	return result
}
