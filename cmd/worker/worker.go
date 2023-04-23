package worker

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/kubeProxy"
	"minik8s/pkg/kubelet"
	"os"
	"os/signal"
	"syscall"
)

var proxyCmd = &cobra.Command{
	Use:   "worker",
	Short: "start the data plane of k8s, this should run on every worker node!",
	Long:  "start the data plane of k8s, this should run on every worker node!",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	// Receive Ctrl-C
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	go kubeProxy.Start_proxy()
	go kubelet.Start_kubelet()

	// Gracefully exit after Ctrl-C
	<-c
	kubelet.ToExit <- true
	<-kubelet.Exited
}

func init() {
}

func Execute() *cobra.Command {
	return proxyCmd
}
