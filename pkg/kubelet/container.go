package kubelet

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"minik8s/pkg/object"
	"minik8s/pkg/util/network"
	"minik8s/pkg/util/weave"
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

func CreateContainer(name string, config *CreateConfig) (string, error) {
	res, err := Client.ContainerCreate(Ctx, &container.Config{
		Image:        config.Image,
		Labels:       config.Labels,
		Entrypoint:   config.Entrypoint,
		Cmd:          config.Cmd,
		Env:          config.Env,
		Volumes:      config.Volumes,
		ExposedPorts: config.ExposedPorts,
	}, &container.HostConfig{
		IpcMode:      config.IpcMode,
		PidMode:      config.PidMode,
		NetworkMode:  config.NetworkMode,
		PortBindings: config.PortBindings,
		Links:        config.Links,
		Binds:        config.Binds,
		VolumesFrom:  config.VolumesFrom,
	}, nil, nil, name)
	return res.ID, err
}

func CreateCommonContainer(pod *object.Pod, myContainer *object.Container) (string, string, error) {
	podFullName := pod.FullName()
	podUuid := pod.Runtime.Uuid

	// Step 1: Prepare for labels
	labels := map[string]string{
		KubernetesPodUIDLabel: podUuid,
	}
	for labelName, labelValue := range pod.Metadata.Labels {
		labels[labelName] = labelValue
	}

	// Step 2: Prepare for CNI
	pauseContainerFullName := pauseContainerFullName(podFullName, podUuid)
	pauseContainerRef := pauseContainerReference(podFullName, podUuid)

	// Step 3: Finally create the container!
	name := ContainerFullName(myContainer.Name, podFullName, podUuid)
	ID, err := CreateContainer(name, &CreateConfig{
		// Config
		Image:      myContainer.Image,
		Labels:     labels,
		Entrypoint: myContainer.Command,
		Cmd:        myContainer.Args,
		Env:        getFormatEnv(myContainer.Env),
		Volumes:    nil,

		// HostConfig
		IpcMode:     container.IpcMode(pauseContainerRef),
		PidMode:     container.PidMode(pauseContainerRef),
		NetworkMode: container.NetworkMode(pauseContainerRef),
		Binds:       getVolumeBinds(pod, myContainer),
		VolumesFrom: []string{pauseContainerFullName},
	})
	return name, ID, err
}

func StartCommonContainer(pod *object.Pod, myContainer *object.Container) {
	// Step 1: Prepare for image
	err := PullImage(myContainer.Image, &PullConfig{
		All: false,
	})
	if err != nil {
		fmt.Printf("Failed to pull image %v! Reason: %v\n", myContainer.Image, err.Error())
	} else {
		fmt.Printf("Image %v pulled!\n", myContainer.Image)
	}

	// Step 2: Create a container
	var fullName, ID string
	fullName, ID, err = CreateCommonContainer(pod, myContainer)
	if err != nil {
		fmt.Printf("Failed to create container %v! Reason: %v\n", fullName, err.Error())
	} else {
		fmt.Printf("Container %v created!\n", fullName)
	}

	// Step 3: Start the container
	err = Client.ContainerStart(Ctx, ID, StartConfig{})
	if err != nil {
		fmt.Printf("Failed to start container %v (ID: %v)! Reason: %v\n", fullName, ID, err.Error())
	} else {
		fmt.Printf("Container %v (ID: %v) started!\n", fullName, ID)
	}
}

func CreatePauseContainer(pod *object.Pod) (string, string, error) {
	podFullName := pod.FullName()
	podUuid := pod.Runtime.Uuid

	// Step 1: Prepare for labels
	labels := map[string]string{
		KubernetesPodUIDLabel: podUuid,
	}
	for labelName, labelValue := range pod.Metadata.Labels {
		labels[labelName] = labelValue
	}

	// Step 2: All the containers share network namespace with pause container
	portBindings := nat.PortMap{}
	portSet := nat.PortSet{}
	for _, c := range pod.Spec.Containers {
		addPortBindings(portBindings, c.Ports)
		addPortSet(portSet, c.Ports)
	}

	// Step 3: Finally create the container!
	name := pauseContainerFullName(podFullName, podUuid)
	ID, err := CreateContainer(name, &CreateConfig{
		// Config
		Image:        pauseImage,
		Labels:       labels,
		Volumes:      nil,
		ExposedPorts: portSet,

		// HostConfig
		IpcMode:      "shareable",
		PortBindings: portBindings,
		Binds:        nil,
	})
	return name, ID, err
}

func StartPauseContainer(pod *object.Pod) {
	// Step 1: Prepare for image
	err := PullImage(pauseImage, &PullConfig{
		All: false,
	})
	if err != nil {
		fmt.Printf("Failed to pull pause image %v! Reason: %v\n", pauseImage, err.Error())
	} else {
		fmt.Printf("Pause image %v pulled!\n", pauseImage)
	}

	// Step 2: Create a container
	var fullName, ID string
	fullName, ID, err = CreatePauseContainer(pod)
	if err != nil {
		fmt.Printf("Failed to create pause container %v! Reason: %v\n", fullName, err.Error())
	} else {
		fmt.Printf("Pause container %v created!\n", fullName)
	}

	// Step 3: Start the container
	err = Client.ContainerStart(Ctx, ID, StartConfig{})
	if err != nil {
		fmt.Printf("Failed to start pause container %v (ID: %v)! Reason: %v\n", fullName, ID, err.Error())
	} else {
		fmt.Printf("Pause container %v (ID: %v) started!\n", fullName, ID)
	}

	// Step 4: Attach to weave subnet
	err = weave.Attach(ID, pod.Runtime.ClusterIp+network.Mask)
	if err != nil {
		fmt.Printf("Failed to attach pause container %v (ID: %v) to subnet! Reason: %v\n", fullName, ID, err.Error())
	} else {
		fmt.Printf("Pause Container %v (ID: %v) attached to subnet!\n", fullName, ID)
	}
}
