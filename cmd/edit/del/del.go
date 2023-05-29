package del

import (
	"github.com/spf13/cobra"
	"minik8s/cmd/edit"
	"minik8s/pkg/client"
	"minik8s/pkg/util/config"
)

var file string
var key string
var tp string

var delCmd = &cobra.Command{
	Use:   "delete",
	Short: "del the items of k8s",
	Long:  "this is the main cmd to delete the items int k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	client.Delete_object(key, tp)
	edit.DelLog(key, tp)
}

func init() {
	delCmd.Flags().StringVarP(&tp, "type", "t", "", "Type of API object(s) to delete")
	delCmd.MarkFlagRequired("type")
	delCmd.Flags().StringVarP(&key, "key", "k", config.EMPTY_FLAG, "Name of API object to delete, refers to all API objects of specified type if not set")
}

func Delete() *cobra.Command {
	return delCmd
}
