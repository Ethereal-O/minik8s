package kubelet

const (
	pauseImage         = "registry.aliyuncs.com/google_containers/pause:3.6"
	pauseContainerName = "PAUSE"
)

const (
	NAMESPACE        = "minik8s"
	JOBNAME          = "resource_usage"
	POD_NAME_PRIFIX  = "pod"
	POD_SUBSYS       = "podResource"
	NODE_NAME_PRIFIX = "node"
	NODE_SUBSYS      = "nodeResource"
)
