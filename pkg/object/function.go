package object

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
