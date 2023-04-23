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

type ResourcesUsage struct {
	CpuPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
}

// Status represents the status of a container.
type Status struct {
	ID             string
	Name           string
	State          State
	CreateTime     time.Time
	StartTime      time.Time
	FinishTime     time.Time
	ExitCode       int
	ImageID        string
	RestartCount   int
	Error          string
	ResourcesUsage ResourcesUsage
	PortBindings   nat.PortMap
}

type ContainerRuntime struct {
	ID      string
	Name    string
	Image   string
	ImageID string
	State   State
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
}

// StartConfig : arguments to start a container
type StartConfig = types.ContainerStartOptions

// ------------------Image------------------

type PullConfig struct {
	All bool
}
