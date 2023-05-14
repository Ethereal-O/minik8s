package apiServer

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/structure"
)

var monitorMap structure.DoubleDirectionMap
var monitorEtcdStopMap map[string]context.CancelFunc
var monitorProducerStopMap map[string]context.CancelFunc

func Init_server() {
	monitorMap.Init()
	monitorEtcdStopMap = make(map[string]context.CancelFunc)
	monitorProducerStopMap = make(map[string]context.CancelFunc)
}

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_server() {
	e := echo.New()
	e.POST("/", basic_post)
	e.PUT("/Pod/:key", pod_put)
	e.GET("/Pod/:key", pod_get)
	e.DELETE("/Pod/:key", pod_delete)
	e.PUT("/ReplicaSet/:key", replicaset_put)
	e.GET("/ReplicaSet/:key", replicaset_get)
	e.DELETE("/ReplicaSet/:key", replicaset_delete)
	e.PUT("/AutoScaler/:key", autoscaler_put)
	e.GET("/AutoScaler/:key", autoscaler_get)
	e.DELETE("/AutoScaler/:key", autoscaler_delete)
	e.PUT("/Node/:key", node_put)
	e.GET("/Node/:key", node_get)
	e.DELETE("/Node/:key", node_delete)
	e.PUT("/Service/:key", service_put)
	e.GET("/Service/:key", service_get)
	e.DELETE("/Service/:key", service_delete)
	e.PUT("/RuntimeService/:key", runtimeService_put)
	e.GET("/RuntimeService/:key", runtimeService_get)
	e.DELETE("RuntimeService/:key", runtimeService_delete)
	e.PUT("/Gateway/:key", gateway_put)
	e.GET("/Gateway/:key", gateway_get)
	e.DELETE("Gateway/:key", gateway_delete)
	e.PUT("/RuntimeGateway/:key", runtimeGateway_put)
	e.GET("/RuntimeGateway/:key", runtimeGateway_get)
	e.DELETE("RuntimeGateway/:key", runtimeGateway_delete)
	e.PUT("/GpuJob/:key", gpujob_put)
	e.GET("/GpuJob/:key", gpujob_get)
	e.DELETE("/GpuJob/:key", gpujob_delete)
	e.PUT("/ServerlessFunctions/:key", serverlessFunctions_put)
	e.GET("/ServerlessFunctions/:key", serverlessFunctions_get)
	e.DELETE("/ServerlessFunctions/:key", serverlessFunctions_delete)
	e.PUT("/Function/:key", function_put)
	e.GET("/Function/:key", function_get)
	e.DELETE("/Function/:key", function_delete)
	e.PUT("/TransFile/:key", transfile_put)
	e.GET("/TransFile/:key", transfile_get)
	go func() { e.Logger.Fatal(e.Start(":8080")) }()

	fmt.Println("API Server start at " + config.APISERVER_URL)

	// Wait until Ctrl-C
	<-ToExit
	Exited <- true
}
