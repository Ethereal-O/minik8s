package worker

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/kubeProxy"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/object"
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

	pod := object.Pod{}
	client.AddPod(pod)

	kubeProxy.Init()
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
