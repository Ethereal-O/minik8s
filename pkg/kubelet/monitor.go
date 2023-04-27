package kubelet

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	namePrefix = "pod"
	subSys     = "pod_resource"
	nameSpace  = "minik8s"
)

var memoryPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: subSys,
		Namespace: nameSpace,
		Name:      fmt.Sprintf("%s:%s", namePrefix, "memory"),
	}, []string{"uuid", "podName"})

var cpuPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: subSys,
		Namespace: nameSpace,
		Name:      fmt.Sprintf("%s:%s", namePrefix, "cpu"),
	}, []string{"uuid", "podName"})

func start_monitor() {
	prometheus.MustRegister(memoryPrecentage, cpuPrecentage)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9080", nil)
}
