package kubeProxy

import (
	"bufio"
	"fmt"
	"minik8s/pkg/object"
	"minik8s/pkg/services"
	"os"
)

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
