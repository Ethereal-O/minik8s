package kubelet

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"minik8s/pkg/object"
	"minik8s/pkg/util/random"
)

var Ctx = context.Background()

var Client = newClient()

func newClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return cli
}

func CreateContainer(name string, config *CreateConfig) {
	_, err := Client.ContainerCreate(Ctx, &container.Config{
		Image:        config.Image,
		Labels:       config.Labels,
		Entrypoint:   config.Entrypoint,
		Cmd:          config.Cmd,
		Env:          config.Env,
		Volumes:      config.Volumes,
		ExposedPorts: config.ExposedPorts,
		Tty:          config.Tty,
	}, &container.HostConfig{
		IpcMode:      config.IpcMode,
		PidMode:      config.PidMode,
		NetworkMode:  config.NetworkMode,
		PortBindings: config.PortBindings,
		Links:        config.Links,
		Binds:        config.Binds,
		VolumesFrom:  config.VolumesFrom,
	}, nil, nil, name)
	if err != nil {
		fmt.Printf("Failed to create container %v! Reason: %v\n", name, err.Error())
	} else {
		fmt.Printf("Container %v created!\n", name)
	}
}

func CreatePod(pod *object.Pod) {
	for _, myContainer := range pod.Spec.Containers {
		ports := nat.PortSet{}
		for _, port := range myContainer.Ports {
			if port.Protocol == "" {
				port.Protocol = "tcp"
			}
			// natPort is a string containing port number and protocol in the format "80/tcp"
			natPort := nat.Port(string(rune(port.ContainerPort)) + "/" + port.Protocol)
			ports[natPort] = struct{}{}
		}

		volumes := map[string]struct{}{}
		for _, volume := range myContainer.VolumeMounts {
			volumes[volume.MountPath] = struct{}{}
		}
		err := PullImage(myContainer.Image, &PullConfig{
			All:     false,
			Verbose: true,
		})
		if err != nil {
			fmt.Printf("Failed to pull image %v! Reason: %v\n", myContainer.Image, err.Error())
		}
		CreateContainer(pod.Metadata.Name+"_"+myContainer.Name+"_"+random.String(6), &CreateConfig{
			Image:      myContainer.Image,
			Entrypoint: myContainer.Command,
			Cmd:        myContainer.Args,
			Volumes:    nil,
			Binds:      toVolumeBinds(pod, &myContainer),
		})
	}
}
