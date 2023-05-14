package kubeProxy

import (
	"bufio"
	"fmt"
	"io"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/services"
	"os"
	"os/exec"
	"strings"
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

func createDir(path string) {
	args := fmt.Sprintf("-r %s %s", services.NGINX_TEMPLATE_FILEPATH, path)
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
	res, err := execCmd("rm", args)
	if err != nil {
		fmt.Println("deleteDir fail")
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
