package apiServer

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/apiServer"
)

var apiCmd = &cobra.Command{
	Use:   "apiServer",
	Short: "start the apiServer of k8s",
	Long:  "start the apiServer of k8s and will always run",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	fmt.Println("apiServer start at 192.168.142.135:8080")
	apiServer.Init_server()
	apiServer.Start_server()
}

func init() {
}

func Api() *cobra.Command {
	return apiCmd
}
