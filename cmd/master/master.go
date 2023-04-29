package master

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/apiServer"
	"minik8s/pkg/controller"
	"minik8s/pkg/scheduler"
	"minik8s/pkg/services"
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
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	apiServer.Init_server()
	go apiServer.Start_server()
	// Wait for API Server to start
	time.Sleep(1 * time.Second)
	go controller.Start_rsController()
	go controller.Start_hpaController()
	go scheduler.Start_scheduler()
	go services.StartServiceManager()
	go services.StartDnsManager()

	// Gracefully exit after Ctrl-C
	<-c
	controller.RSToExit <- true
	controller.HpaToExit <- true
	scheduler.ToExit <- true
	services.ToExit <- true
	<-controller.RSExited
	<-controller.HpaExited
	<-scheduler.Exited
	<-services.Exited

	// Wait for other components to exit
	apiServer.ToExit <- true
	<-apiServer.Exited
}

func init() {
}

func Execute() *cobra.Command {
	return apiCmd
}
