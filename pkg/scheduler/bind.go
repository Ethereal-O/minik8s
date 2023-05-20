package scheduler

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"minik8s/pkg/client"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/resource"
	"time"
)

func BindPod(pod *object.Pod, policy SchedulePolicy) bool {
	optional_nodes := getOptionalNodes(pod)
	pod.Runtime.Bind = policy.selectNode(pod, optional_nodes)
	if pod.Runtime.Bind != "" {
		pod.Runtime.Status = config.BOUND_STATUS
		client.AddPod(*pod)
		return true
	} else {
		return false
	}
}

type SchedulePolicy interface {
	selectNode(pod *object.Pod, nodes []object.Node) string
}

func getOptionalNodes(pod *object.Pod) []object.Node {
	nodes := client.GetActiveNodes()

	var optional_nodes []object.Node
	for _, node := range nodes {
		cpu := prometheus_query("cpu", node.Runtime.Uuid)
		mem := prometheus_query("memory", node.Runtime.Uuid)
		for _, container := range pod.Spec.Containers {
			cpu -= float64(resource.ConvertCpuToBytes(container.Limits.Cpu))
			mem -= float64(resource.ConvertMemoryToBytes(container.Limits.Memory))
		}
		if cpu > 0 && mem > 0 {
			optional_nodes = append(optional_nodes, node)
		}
	}

	return optional_nodes
}

func prometheus_query(tp string, uuid string) float64 {

	// Step 1: Form the query with type(cpu/memory) and uuid
	publicPrix := fmt.Sprintf("%s_%s_%s",
		kubelet.NAMESPACE, kubelet.NODE_SUBSYS, kubelet.NODE_NAME_PRIFIX)
	query := fmt.Sprintf("%s:%s{job=\"%s\",uuid=\"%s\"}",
		publicPrix, tp, kubelet.JOBNAME, uuid)

	// Step 2: Query for the target metric
	prometheusClient, err := api.NewClient(api.Config{
		Address: config.PROMETHEUS_URL,
	})
	if err != nil {
		fmt.Println("[Prometheus Client error]", err.Error())
		panic(err)
	}
	v1api := v1.NewAPI(prometheusClient)
	result, _, err := v1api.Query(context.Background(), query, time.Now())
	if err != nil {
		fmt.Println("[Query error]", err.Error())
		panic(err)
	}

	// Step 3: Print the result (for debug) and return
	// Because uuid is unique,there is only one result for each uuid and each type (cpu/memory)
	// But there can be 0 result since that the node is created but hasn't upload
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) != 0 {
			s := vector[0]
			//fmt.Printf("nodeName=%q, uuid=%q, value=%v\n",
			//	s.Metric["nodeName"], s.Metric["uuid"], s.Value)
			return float64(s.Value)
		}
	}
	return 0.0
}
