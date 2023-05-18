package controller

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/util/config"
	"time"
)

func judge(tp string, strategy string, limit int, uuidList []string) bool {
	fmt.Printf("[strategy]%s\n", strategy)
	if strategy == "" {
		return true
	} else if strategy == "max" {
		var max float64 = 0.0
		for _, uuid := range uuidList {
			tmp := prometheus_query(tp, uuid)
			if tmp > max {
				max = tmp
			}
		}
		fmt.Printf("[%s](%s) actual: %.3f limit :%.3f\n", strategy, tp, max, float64(limit))
		return max <= float64(limit)
	} else if strategy == "average" {
		var sum float64 = 0
		for _, uuid := range uuidList {
			sum += prometheus_query(tp, uuid)
		}
		average := sum / (float64(len(uuidList)))
		fmt.Printf("[%s](%s) actual: %.3f limit :%.3f\n", strategy, tp, average, float64(limit))
		return average <= float64(limit)
	} else {
		fmt.Println("Invalid strategy!")
		return false
	}
}

func prometheus_query(tp string, uuid string) float64 {

	// Step 1: Form the query with type(cpu/memory) and uuid
	publicPrix := fmt.Sprintf("%s_%s_%s",
		kubelet.NAMESPACE, kubelet.POD_SUBSYS, kubelet.POD_NAME_PRIFIX)
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
	// But there can be 0 result since that the pod is created but hasn't upload
	if result.Type() == model.ValVector {
		vector := result.(model.Vector)
		if len(vector) != 0 {
			s := vector[0]
			//fmt.Printf("podName=%q, uuid=%q, value=%v\n",
			//	s.Metric["podName"], s.Metric["uuid"], s.Value)
			return float64(s.Value)
		}
	}
	return 0.0
}
