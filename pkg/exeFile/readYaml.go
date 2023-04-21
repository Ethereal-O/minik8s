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
	} else if strings.Contains(string(yamlFile), "kind: Replicaset") {
		return parseRs(yamlFile)
	} else if strings.Contains(string(yamlFile), "kind: Node") {
		return parseNode(yamlFile)
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
	return string(inf), key, "Replicaset"
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
