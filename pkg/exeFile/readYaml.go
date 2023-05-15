package exeFile

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/pkg/fileServer"
	"minik8s/pkg/object"
	"os"
	"strings"
)

func ReadRequest(file string) map[string]string {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	params := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, &params)
	return params
}

func ReadYaml(file string) (string, string, string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	if strings.Contains(string(yamlFile), "kind: Pod") {
		return parsePod(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: ReplicaSet") {
		return parseRs(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: AutoScaler") {
		return parseAutoScaler(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: Node") {
		return parseNode(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: Service") {
		return parseService(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: Gateway") {
		return parseGateway(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: GpuJob") {
		return parseGpuJob(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: ServerlessFunctions") {
		return parseServerlessFunctions(yamlFile)
	} else {
		return "", "", ""
	}
}

func parsePod(yamlFile []byte) (string, string, string) {
	var conf object.Pod
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "Pod"
}

func parseService(yamlFile []byte) (string, string, string) {
	var conf object.Service
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "Service"
}

func parseGateway(yamlFile []byte) (string, string, string) {
	var conf object.Gateway
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "Gateway"
}

func parseRs(yamlFile []byte) (string, string, string) {
	var conf object.ReplicaSet
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "ReplicaSet"
}

func parseAutoScaler(yamlFile []byte) (string, string, string) {
	var conf object.AutoScaler
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "AutoScaler"
}

func parseNode(yamlFile []byte) (string, string, string) {
	var conf object.Node
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	key = conf.Metadata.Name
	inf, err = json.Marshal(&conf)
	return string(inf), key, "Node"
}

func parseGpuJob(yamlFile []byte) (string, string, string) {
	var conf object.GpuJob
	var inf []byte
	var key string
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}

	// the yaml only contains the filepath,the actual file data should be transmitted
	key = conf.Metadata.Name
	path := conf.Spec.Path
	dir, err := os.Getwd()
	conf.Spec.Path = dir + path
	fileServer.UploadFile(dir+path, key, "GpuFile")
	inf, err = json.Marshal(&conf)
	return string(inf), key, "GpuJob"
}

func parseServerlessFunctions(yamlFile []byte) (string, string, string) {
	var conf object.ServerlessFunctions
	err := yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}

	// the yaml only contains the filepath,the actual file data should be transmitted
	key := conf.Metadata.Name
	dir, err := os.Getwd()
	path := conf.Spec.Path
	fileServer.UploadFile(dir+path, key, "FuncFile")

	inf, _ := json.Marshal(conf)
	return string(inf), key, "ServerlessFunctions"

}
