package watch

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/messging"
)

var key string
var prix bool

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch the items of k8s",
	Long:  "this is the main cmd to watch the items int k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	//etcd.Watch_etcd(key, prix)
	messging.Watch(key, prix)
}

func init() {
	watchCmd.Flags().StringVarP(&key, "key", "k", "", "")
	watchCmd.Flags().BoolVarP(&prix, "prix", "p", false, "")
}

func Watch() *cobra.Command {
	return watchCmd
}
