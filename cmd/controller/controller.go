package controller

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/controller"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "start the controller of k8s",
	Long:  "start the controller of k8s and will always run",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	fmt.Println("controller start ")
	controller.Start_rsController()
}

func init() {
}

func StartController() *cobra.Command {
	return controllerCmd
}
