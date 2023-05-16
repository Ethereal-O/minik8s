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
	e.POST("/", forwardRequest)
	go func() { e.Logger.Fatal(e.Start(":8081")) }()
	go FlowControl()

	fmt.Println("Function proxy start at " + config.FUNCTION_PROXY_URL)

	// Wait until Ctrl-C
	<-ToExit
	Exited <- true
}
