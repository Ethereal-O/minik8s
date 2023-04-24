package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io"
	"minik8s/pkg/object"
	"minik8s/pkg/util/network"
	"strconv"
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

func pauseContainerFullName(podFullName string, podUuid string) string {
	return ContainerFullName(pauseContainerName, podFullName, podUuid)
}

func pauseContainerReference(podFullName string, podUuid string) string {
	return "container:" + pauseContainerFullName(podFullName, podUuid)
}

func ContainerFullName(containerName, podFullName, podUuid string) string {
	return podFullName + "_" + podUuid + "_" + containerName
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
