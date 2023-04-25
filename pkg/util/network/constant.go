package network

const (
	PodIPGeneratorURL        = "/ip_generator/pod"
	PodIPGeneratorInitIP     = "10.10.1.1"
	ServiceIPGeneratorURL    = "/ip_generator/service"
	ServiceIPGeneratorInitIP = "10.10.20.1"
	NodeIPGeneratorURL       = "/ip_generator/node"
	NodeIPGeneratorInitIP    = "10.10.0.1"
)

// Mask : e.g. 192.168.1.1/16
const Mask = "/16"
