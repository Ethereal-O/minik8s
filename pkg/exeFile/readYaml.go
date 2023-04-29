package exeFile

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/pkg/object"
	"strings"
)

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
