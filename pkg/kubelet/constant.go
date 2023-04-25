package kubelet

const (
	pauseImage                   = "registry.aliyuncs.com/google_containers/pause:3.6"
	pauseContainerName           = "PAUSE"
	KubernetesPodNameLabel       = "io.kubernetes.pod.name"
	KubernetesPodUIDLabel        = "io.kubernetes.pod.uid"
	KubernetesReplicaSetUIDLabel = "io.kubernetes.rs.uid"
	KubernetesContainerNameLabel = "io.kubernetes.container.name"
)
