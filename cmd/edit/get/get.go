package get

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
)

var key string
var tp string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get the items of k8s",
	Long:  "this is the main cmd to get the items int k8s",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	res := client.Get_object(key, tp)
	if tp == config.POD_TYPE {
		fmt.Println("Name\tUuid\tBelong")
		for _, pod := range res {
			var podObject object.Pod
			json.Unmarshal([]byte(pod), &podObject)
			if podObject.Runtime.Status == config.RUNNING_STATUS {
				fmt.Println(podObject.Metadata.Name, podObject.Runtime.Uuid, podObject.Runtime.Belong)
			}
		}
	}
	if tp == config.REPLICASET_TYPE {
		fmt.Println("Name\tUuid\tReplicas")
		for _, rs := range res {
			var rsObject object.ReplicaSet
			json.Unmarshal([]byte(rs), &rsObject)
			if rsObject.Runtime.Status == config.RUNNING_STATUS {
				fmt.Println(rsObject.Metadata.Name, rsObject.Runtime.Uuid, rsObject.Spec.Replicas)
			}
		}
	}
}

func init() {
	getCmd.Flags().StringVarP(&tp, "type", "t", "", "")
	getCmd.MarkFlagRequired("type")
	getCmd.Flags().StringVarP(&key, "key", "k", config.EMPTY_FLAG, "")
}

func Get() *cobra.Command {
	return getCmd
}
