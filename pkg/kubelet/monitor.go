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
	NAME_PRIFIX = "pod"
	SUBSYS      = "podResource"
	NAMESPACE   = "minik8s"
	JOBNAME     = "resource_usage"
)

var memoryPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", NAME_PRIFIX, "memory"),
	}, []string{"uuid", "podName"})

var cpuPrecentage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Subsystem: SUBSYS,
		Namespace: NAMESPACE,
		Name:      fmt.Sprintf("%s:%s", NAME_PRIFIX, "cpu"),
	}, []string{"uuid", "podName"})

func start_monitor() {
	prometheus.MustRegister(memoryPrecentage, cpuPrecentage)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9080", nil)
}
