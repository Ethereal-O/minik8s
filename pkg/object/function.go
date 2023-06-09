package object

import (
	"minik8s/pkg/util/config"
	"strconv"
)

// --- Pods and Rs/Ds ---

func GetPodsOfRS(rs *ReplicaSet, activePods []Pod) ([]Pod, int) {
	actualNum := 0
	var rspodList []Pod
	for _, pod := range activePods {
		if pod.Runtime.Belong == rs.Metadata.Name {
			actualNum++
			rspodList = append(rspodList, pod)
		}
	}
	return rspodList, actualNum
}

func RSPodFullName(rs *ReplicaSet, pod *Pod) string {
	return rs.Metadata.Name + "_" + pod.Runtime.Uuid
}

func GetPodsOfDS(ds *DaemonSet, activePods []Pod) ([]Pod, int) {
	actualNum := 0
	var dspodList []Pod
	for _, pod := range activePods {
		if pod.Runtime.Belong == ds.Metadata.Name {
			actualNum++
			dspodList = append(dspodList, pod)
		}
	}
	return dspodList, actualNum
}

func DSPodFullName(ds *DaemonSet, node *Node) string {
	return ds.Metadata.Name + "_" + node.Metadata.Name
}

func SerializePodList(podList []Pod) string {
	serialized := ""
	for idx, pod := range podList {
		serialized += pod.Metadata.Name
		if idx < len(podList)-1 {
			serialized += ", "
		}
	}
	return serialized
}

// --- Service ---

func SerializeSelectorList(selectorList map[string]string) string {
	serialized := ""
	i := 0
	for k, v := range selectorList {
		serialized += k + ":" + v
		if i < len(selectorList)-1 {
			serialized += ", "
		}
		i++
	}
	return serialized
}

func SerializeEndPortsList(servicePortList []ServicePort, tp string) string {
	serialized := ""
	for idx, port := range servicePortList {
		if tp == config.SERVICE_TYPE_NODEPORT {
			serialized += port.NodePort + ":" + port.Port + "->" + port.TargetPort + "(" + port.Protocol + ")"
		} else {
			serialized += port.Port + "->" + port.TargetPort + "(" + port.Protocol + ")"
		}
		if idx < len(servicePortList)-1 {
			serialized += ", "
		}
	}
	return serialized
}

func SerializeEndPointsList(podList []Pod) string {
	serialized := ""
	for idx, pod := range podList {
		serialized += pod.Metadata.Name + ":" + pod.Runtime.ClusterIp
		if idx < len(podList)-1 {
			serialized += ", "
		}
	}
	return serialized
}

func SerializePathList(pathList []Path) string {
	serialized := ""
	for idx, path := range pathList {
		serialized += path.Name + "->" + path.Service
		if idx < len(pathList)-1 {
			serialized += ", "
		}
	}
	return serialized
}

// --- VirtualService ---

func SerializeVirtualSelectorList(selectorList []VirtualServiceSelector) string {
	serialized := ""
	i := 0
	for _, selector := range selectorList {
		serialized += "("
		j := 0
		for k, v := range selector.MatchLabels {
			serialized += k + ":" + v
			if j < len(selector.MatchLabels)-1 {
				serialized += ", "
			}
			j++
		}
		serialized += "):" + strconv.Itoa(selector.Weight)
		if i < len(selectorList)-1 {
			serialized += ", "
		}
		i++
	}
	return serialized
}

// --- GpuJob  ---

func GpuJobPodFullName(job GpuJob) string {
	return config.GPU_JOB_NAME + "_pod_" + job.Metadata.Name
}

// --- Serverless  ---

func FaasRsFullName(functions ServerlessFunctions) string {
	return config.FUNC_NAME + "_rs_" + functions.Metadata.Name
}

func FaasServiceFullName(functions ServerlessFunctions) string {
	return config.FUNC_NAME + "_service_" + functions.Metadata.Name
}

func FaasPodFullName(functions ServerlessFunctions) string {
	return config.FUNC_NAME + "_pod_" + functions.Metadata.Name
}
