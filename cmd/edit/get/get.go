package get

import (
	"encoding/json"
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"strconv"
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
	fmt.Printf("Type: %v\n", tp)
	if tp == config.SERVICE_TYPE {
		tp = config.RUNTIMESERVICE_TYPE
	}
	if tp == config.GATEWAY_TYPE {
		tp = config.RUNTIMEGATEWAY_TYPE
	}
	res := client.Get_object(key, tp)
	if key != config.EMPTY_FLAG {
		fmt.Printf("Key: %v\n", key)
	}
	if tp == config.POD_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "Belong", "Bind", "ClusterIP")
		for _, pod := range res {
			var podObject object.Pod
			json.Unmarshal([]byte(pod), &podObject)
			rows := make([]map[string]string, 0)
			if podObject.Runtime.Status != config.EXIT_STATUS {
				row := make(map[string]string)
				row["Name"] = podObject.Metadata.Name
				row["Uuid"] = podObject.Runtime.Uuid
				row["Status"] = podObject.Runtime.Status
				row["Belong"] = podObject.Runtime.Belong
				row["Bind"] = podObject.Runtime.Bind
				row["ClusterIP"] = podObject.Runtime.ClusterIp
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}

	if tp == config.REPLICASET_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "Replicas", "ActualReplicas", "Pods")
		for _, rs := range res {
			var rsObject object.ReplicaSet
			json.Unmarshal([]byte(rs), &rsObject)
			rows := make([]map[string]string, 0)
			if rsObject.Runtime.Status != config.EXIT_STATUS {
				rspodList, actualNum := object.GetPodsOfRS(&rsObject, client.GetActivePods())

				row := make(map[string]string)
				row["Name"] = rsObject.Metadata.Name
				row["Uuid"] = rsObject.Runtime.Uuid
				row["Status"] = rsObject.Runtime.Status
				row["Replicas"] = strconv.Itoa(rsObject.Spec.Replicas)
				row["ActualReplicas"] = strconv.Itoa(actualNum)
				row["Pods"] = object.SerializePodList(rspodList)
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}
	if tp == config.AUTOSCALER_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "MinReplicas", "MaxReplicas", "ActualReplicas")
		for _, hpa := range res {
			var hpaObject object.AutoScaler
			json.Unmarshal([]byte(hpa), &hpaObject)
			rows := make([]map[string]string, 0)
			if hpaObject.Runtime.Status != config.EXIT_STATUS {
				rsList := client.GetReplicaSetByKey(hpaObject.Spec.ScaleTargetRef.Name)
				tarRs := rsList[0]
				_, actualNum := object.GetPodsOfRS(&tarRs, client.GetActivePods())

				row := make(map[string]string)
				row["Name"] = hpaObject.Metadata.Name
				row["Uuid"] = hpaObject.Runtime.Uuid
				row["Status"] = hpaObject.Runtime.Status
				row["MinReplicas"] = strconv.Itoa(hpaObject.Spec.MinReplicas)
				row["MaxReplicas"] = strconv.Itoa(hpaObject.Spec.MaxReplicas)
				row["ActualReplicas"] = strconv.Itoa(actualNum)
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}
	if tp == config.NODE_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "PublicIP", "ClusterIP")
		for _, node := range res {
			var nodeObject object.Node
			json.Unmarshal([]byte(node), &nodeObject)
			rows := make([]map[string]string, 0)
			if nodeObject.Runtime.Status != config.EXIT_STATUS {
				row := make(map[string]string)
				row["Name"] = nodeObject.Metadata.Name
				row["Uuid"] = nodeObject.Runtime.Uuid
				row["Status"] = nodeObject.Runtime.Status
				row["PublicIP"] = nodeObject.Spec.Ip
				row["ClusterIP"] = nodeObject.Runtime.ClusterIp
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}
	if tp == config.RUNTIMESERVICE_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "Selector", "Type", "IP", "Port", "Endpoints")
		for _, service := range res {
			var runtimeServiceObject object.RuntimeService
			json.Unmarshal([]byte(service), &runtimeServiceObject)
			rows := make([]map[string]string, 0)
			if runtimeServiceObject.Service.Runtime.Status != config.EXIT_STATUS {
				row := make(map[string]string)
				row["Name"] = runtimeServiceObject.Service.Metadata.Name
				row["Uuid"] = runtimeServiceObject.Service.Runtime.Uuid
				row["Status"] = runtimeServiceObject.Service.Runtime.Status
				row["Selector"] = object.SerializeSelectorList(runtimeServiceObject.Service.Spec.Selector)
				row["Type"] = runtimeServiceObject.Service.Spec.Type
				row["IP"] = runtimeServiceObject.Service.Runtime.ClusterIp
				row["Port"] = object.SerializeEndPortsList(runtimeServiceObject.Service.Spec.Ports, runtimeServiceObject.Service.Spec.Type)
				row["Endpoints"] = object.SerializeEndPointsList(runtimeServiceObject.Pods)
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}
	if tp == config.RUNTIMEGATEWAY_TYPE {
		table, _ := gotable.Create("Name", "Uuid", "Status", "Host", "Path")
		for _, gateway := range res {
			var runtimeGatewayObject object.RuntimeGateway
			json.Unmarshal([]byte(gateway), &runtimeGatewayObject)
			rows := make([]map[string]string, 0)
			if runtimeGatewayObject.Gateway.Runtime.Status != config.EXIT_STATUS {
				row := make(map[string]string)
				row["Name"] = runtimeGatewayObject.Gateway.Metadata.Name
				row["Uuid"] = runtimeGatewayObject.Gateway.Runtime.Uuid
				row["Status"] = runtimeGatewayObject.Gateway.Runtime.Status
				row["Host"] = runtimeGatewayObject.Gateway.Spec.Host
				row["Path"] = object.SerializePathList(runtimeGatewayObject.Gateway.Spec.Paths)
				rows = append(rows, row)
			}
			table.AddRows(rows)
		}
		fmt.Println(table)
	}
}

func init() {
	getCmd.Flags().StringVarP(&tp, "type", "t", "", "Type of API object(s) to inspect")
	getCmd.MarkFlagRequired("type")
	getCmd.Flags().StringVarP(&key, "key", "k", config.EMPTY_FLAG, "Name of API object to inspect, refers to all API objects of specified type if not set")
}

func Get() *cobra.Command {
	return getCmd
}
