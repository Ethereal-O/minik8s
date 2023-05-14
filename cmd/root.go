package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/cmd/edit/apply"
	"minik8s/cmd/edit/del"
	"minik8s/cmd/edit/get"
	"minik8s/cmd/edit/watch"
	"minik8s/cmd/master"
	"minik8s/cmd/request"
	"minik8s/cmd/worker"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "kubelet",
	Short: "manage the items of k8s",
	Long:  "this is the main cmd to controll the items int k8s",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("this is kubelet")
	},
}

func init() {
	rootCmd.AddCommand(apply.Apply(), del.Delete(), get.Get(), watch.Watch(),
		master.Execute(), worker.Execute(), request.Request())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}
