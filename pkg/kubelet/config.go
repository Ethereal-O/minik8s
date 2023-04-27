package kubelet

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"time"
)

// ------------------Container------------------

type State string

const (
	StateCreated State = "CREATED"
	StateRunning State = "RUNNING"
	StateExited  State = "EXITED"
	StateUnknown State = "UNKNOWN"
)

// Status represents the status of a container.
type Status struct {
	ID           string
	Name         string
	State        State
	CreatedAt    time.Time
	StartedAt    time.Time
	FinishedAt   time.Time
	ExitCode     int
	ImageID      string
	RestartCount int
	Error        string
	PortBindings nat.PortMap
	CpuPercent   float64
	MemPercent   float64
}

// CreateConfig : arguments to create a container
type CreateConfig struct {
	// Config
	Image        string
	Labels       map[string]string
	Entrypoint   []string
	Cmd          []string
	Env          []string
	Volumes      map[string]struct{}
	ExposedPorts nat.PortSet // Exposed ports of the container

	// HostConfig
	IpcMode      container.IpcMode     // IPC namespace of the container
	PidMode      container.PidMode     // PID namespace of the container
	NetworkMode  container.NetworkMode // Network mode of the container (e.g., --network=container:nginx)
	PortBindings nat.PortMap           // Port mapping between exposed ports and the host ports
	Links        []string              // List of links (name:alias)
	Binds        []string              // List of volume bindings of the container
	VolumesFrom  []string              // List of volumes to take from other containers
	Memory       int64                 // Memory Limit of the container(in bytes)
	NanoCPUs     int64                 // CPU Limit of the container(in units of 10<sup>-9</sup> CPUs)
}

// StartConfig : arguments to start a container
type StartConfig = types.ContainerStartOptions

// InspectInfo : results of inspecting a container
type InspectInfo = types.ContainerJSON

// StopConfig : arguments to stop a container
type StopConfig = types.ContainerRemoveOptions

// ------------------Image------------------
