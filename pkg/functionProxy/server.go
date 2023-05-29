package functionProxy

import (
	"fmt"
	"github.com/labstack/echo"
	"minik8s/pkg/util/config"
)

var Exited = make(chan bool)
var ToExit = make(chan bool)

func Start_proxy() {
	e := echo.New()
	e.POST("/simpleRequest", forwardRequest)
	e.POST("/workflow", doWorkflow)
	go func() { e.Logger.Fatal(e.Start(":8081")) }()
	go FlowControl()
	cache.Init(10)

	fmt.Println("Function proxy start at " + config.FUNCTION_PROXY_URL)

	// Wait until Ctrl-C
	<-ToExit
	Exited <- true
}
