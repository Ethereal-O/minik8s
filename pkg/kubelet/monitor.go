package kubelet

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var podMemoryPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: POD_SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", POD_NAME_PRIFIX, "memory"),
	}, []string{"uuid", "podName"})

var podCpuPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: POD_SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", POD_NAME_PRIFIX, "cpu"),
	}, []string{"uuid", "podName"})

var nodeAvailableMemory = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: NODE_SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", NODE_NAME_PRIFIX, "memory"),
	}, []string{"uuid", "nodeName"})

var nodeAvailableCpu = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: NODE_SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", NODE_NAME_PRIFIX, "cpu"),
	}, []string{"uuid", "nodeName"})

func start_monitor() {
	prometheus.MustRegister(podMemoryPrecentage, podCpuPrecentage, nodeAvailableMemory, nodeAvailableCpu)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9080", nil)
}
