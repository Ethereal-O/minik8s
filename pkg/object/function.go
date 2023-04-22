package object

func (pod *Pod) FullName() string {
	return pod.Metadata.Name + "_" + pod.Metadata.Namespace
}
