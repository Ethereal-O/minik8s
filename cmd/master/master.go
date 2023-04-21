package master

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/apiServer"
	"minik8s/pkg/controller"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var apiCmd = &cobra.Command{
	Use:   "master",
	Short: "start the control plane of k8s, this should only run on master node!",
	Long:  "start the control plane of k8s, this should only run on master node!",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	// Receive Ctrl-C
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	apiServer.Init_server()
	go apiServer.Start_server()
	// Wait for API Server to start
	time.Sleep(1 * time.Second)
	go controller.Start_rsController()

	// Gracefully exit after Ctrl-C
	<-c
	apiServer.ToExit <- true
	controller.RSToExit <- true
	<-apiServer.Exited
	<-controller.RSExited
}

func init() {
}

func Execute() *cobra.Command {
	return apiCmd
}
