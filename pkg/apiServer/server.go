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
	e.PUT("/Replicaset/:key", replicaset_put)
	e.GET("/Replicaset/:key", replicaset_get)
	e.DELETE("/Replicaset/:key", replicaset_delete)
	e.PUT("/Node/:key", node_put)
	e.GET("/Node/:key", node_get)
	e.DELETE("/Node/:key", node_delete)
	e.PUT("/Service/:key", service_put)
	e.GET("/Service/:key", service_get)
	e.DELETE("/Service/:key", service_delete)
	e.PUT("/ServiceStatus/:key", serviceStatus_put)
	e.GET("/ServiceStatus/:key", serviceStatus_get)
	e.DELETE("/ServiceStatus/:key", serviceStatus_delete)
	go func() { e.Logger.Fatal(e.Start(":8080")) }()
	fmt.Println("API Server start at " + config.APISERVER_URL)

	// Wait until Ctrl-C
	<-ToExit
	Exited <- true
}
