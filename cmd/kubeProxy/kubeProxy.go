package kubeProxy

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/kubeProxy"
)

var proxyCmd = &cobra.Command{
	Use:   "kubeProxy",
	Short: "start the kubeProxy of k8s",
	Long:  "start the kubeProxy of k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	fmt.Println("kubeProxy running")
	kubeProxy.Init()
}

func init() {
}

func Proxy() *cobra.Command {
	return proxyCmd
}
