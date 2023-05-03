package object

import "minik8s/pkg/util/config"

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
