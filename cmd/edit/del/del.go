package del

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
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
	if file != "" {
		_, keyy, tp := exeFile.ReadYaml(file)
		client.Delete_object(keyy, tp)
	} else {
		client.Delete_object(key, tp)
	}
}

func init() {
	delCmd.Flags().StringVarP(&file, "file", "f", "", "")
	delCmd.Flags().StringVarP(&tp, "type", "t", "", "")
	delCmd.Flags().StringVarP(&key, "key", "k", config.EMPTY_FLAG, "")
}

func Delete() *cobra.Command {
	return delCmd
}
