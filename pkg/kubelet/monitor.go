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
	http.HandleFunc("/hello", Hello)
	http.ListenAndServe(":9080", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle hello") // 服务端打印输出
	fmt.Fprintf(w, "hello GoLangWEB")
}

//func start_monitor() {
//	prometheus.MustRegister(memoryPrecentage, cpuPrecentage)
//	http.Handle("/metrics", promhttp.Handler())
//	http.HandleFunc("/hello", Hello)
//
//	// 创建一个信号通道
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//	// 创建一个 HTTP 服务器
//	server := &http.Server{Addr: ":9080"}
//
//	// 启动 HTTP 服务器并监听信号
//	go func() {
//		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			fmt.Println("listen: %s\n", err)
//		}
//	}()
//	fmt.Println("Server started")
//
//	// 等待信号
//	<-sigs
//
//	// 收到信号后，执行一些清理操作
//	fmt.Println("Shutting down server...")
//	if err := server.Shutdown(context.Background()); err != nil {
//		fmt.Printf("Server shutdown failed: %s\n", err)
//	}
//	fmt.Println("Server stopped")
//}
