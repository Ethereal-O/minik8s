package counter

const (
	UUIDCounterURL       = "/counter/uuidCounter"
	UUIDCounterInitCount = "10000"

	MonitorCounterURL       = "/counter/monitor"
	MonitorCounterInitCount = "10000"

	RRPolicyCounterURL       = "/counter/rr_policy"
	RRPolicyCounterInitCount = "0"

	PodIPCounterURL    = "/counter/pod_ip"
	PodIPCounterInitIP = "10.10.1.1"

	ServiceIPCounterURL    = "/counter/service_ip"
	ServiceIPCounterInitIP = "100.10.1.1"

	NodeIPCounterURL    = "/counter/node_ip"
	NodeIPCounterInitIP = "10.10.0.1"

	NodePortCounterURL      = "/counter/node_port"
	NodePortCounterInitPort = "30000"
)
