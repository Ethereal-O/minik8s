package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"io"
	"math"
	"minik8s/pkg/object"
	"minik8s/pkg/util/network"
	"strconv"
	"strings"
)

// ------------------Container------------------

// getVolumeBinds returns the binds of volumes
func getVolumeBinds(pod *object.Pod, target *object.Container) []string {
	// Get volume devices, create a map
	// Mapping volume name to its source
	volumes := make(map[string]string)
	// Now we only support HostPath
	for _, volume := range pod.Spec.Volumes {
		if volume.Type == "hostPath" {
			volumes[volume.Name] = volume.Path
		}
	}

	var volumeBinds []string
	for _, volumeMount := range target.VolumeMounts {
		volumeName := volumeMount.Name
		// If the specified volume device is existent, and is hostPath(we only support this type temporarily)
		if device, exists := volumes[volumeName]; exists {
			// Volume bind rule: $(host path):$(container path)
			mountRule := device + ":" + volumeMount.MountPath
			volumeBinds = append(volumeBinds, mountRule)
		}
	}
	return volumeBinds
}

// getFormatEnv changes container.Env to formatted form, like "PATH=/usr/bin"
func getFormatEnv(containerEnv []object.EnvVar) []string {
	var formatEnv []string
	for _, env := range containerEnv {
		formatEnv = append(formatEnv, env.Name+"="+env.Value)
	}
	return formatEnv
}

func addPortBindings(portBindings nat.PortMap, ports []object.Port) {
	for _, port := range ports {
		// Protocol not assigned, default is tcp
		if port.Protocol == "" {
			port.Protocol = "tcp"
		}
		// HostIP not assigned, default is localhost (127.0.0.1)
		if port.HostIP == "" {
			port.HostIP = "127.0.0.1"
		}
		// HostPort not assigned, default is random available port
		if port.HostPort == 0 {
			randomPort, err := network.GetAvailablePort()
			if err != nil {

			}
			fmt.Printf("Using random available port %d\n", randomPort)
			port.HostPort = randomPort
		}

		// Finally bind them!
		containerPort, err := nat.NewPort(port.Protocol, strconv.Itoa(port.ContainerPort))
		if err != nil {

		}
		portBindings[containerPort] = []nat.PortBinding{{
			HostIP:   port.HostIP,
			HostPort: strconv.Itoa(port.HostPort),
		}}
	}
}

func addPortSet(portSet nat.PortSet, ports []object.Port) {
	for _, port := range ports {
		// Protocol not assigned, default is tcp
		if port.Protocol == "" {
			port.Protocol = "tcp"
		}
		portSet[nat.Port(strconv.Itoa(port.ContainerPort)+"/"+port.Protocol)] = struct{}{}
	}
}

func pauseContainerFullName(podName string, podUuid string) string {
	return ContainerFullName(pauseContainerName, podName, podUuid)
}

func pauseContainerReference(podName string, podUuid string) string {
	return "container:" + pauseContainerFullName(podName, podUuid)
}

func ContainerFullName(containerName, podName, podUuid string) string {
	return podName + "_" + podUuid + "_" + containerName
}

// -------------container resource-------------
func convertMemoryToBytes(memoryStr string) int64 {
	var bytes int64 = 1024 * 1024 * 200
	memoryStr = strings.TrimSpace(memoryStr)
	if memoryStr == "" {
		return bytes //default 200 MB
	}
	memoryStr = strings.ToLower(memoryStr)
	if memoryStr[len(memoryStr)-2:] == "gi" {
		val, err := strconv.ParseFloat(memoryStr[:len(memoryStr)-2], 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = int64(math.Round(val * math.Pow(1024, 3)))
	} else if memoryStr[len(memoryStr)-2:] == "mi" {
		val, err := strconv.ParseFloat(memoryStr[:len(memoryStr)-2], 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = int64(math.Round(val * math.Pow(1024, 2)))
	} else if memoryStr[len(memoryStr)-1:] == "k" {
		val, err := strconv.ParseInt(memoryStr[:len(memoryStr)-1], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = val * 1024
	} else if memoryStr[len(memoryStr)-1:] == "b" {
		val, err := strconv.ParseInt(memoryStr[:len(memoryStr)-1], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return bytes
		}
		bytes = val
	} else {
		return bytes
	}
	return bytes
}

func convertCpuToBytes(cpuStr string) int64 {
	cpuStr = strings.TrimSpace(cpuStr)
	if len(cpuStr) == 0 {
		return 1 * 1e9 //default 1 core
	} else {
		cpuLimit, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			fmt.Println(err.Error())
			return -1
		}
		NanoCPU := (int64)(cpuLimit * 1e9)
		return NanoCPU
	}
}

func calculateCPUPercent(stats types.StatsJSON) float64 {
	cpuPercent := 0.0
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func calculateMemPercent(stats types.StatsJSON) float64 {
	memPercent := float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100.0
	return memPercent
}

// ------------------Image------------------

func waitForPullComplete(events io.ReadCloser) {
	d := json.NewDecoder(events)

	type Event struct {
		Status         string `json:"status"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
		Error string `json:"error"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}
}
