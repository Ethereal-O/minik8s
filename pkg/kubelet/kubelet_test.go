package kubelet

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"minik8s/pkg/object"
	"os"
	"testing"
)

func TestCreatePod(t *testing.T) {
	yamlFile, err := os.ReadFile("../../yaml/pod1.yaml")
	assert.Nil(t, err)
	var conf object.Pod
	err = yaml.Unmarshal(yamlFile, &conf)
	assert.Nil(t, err)
	StartPod(&conf)
}
