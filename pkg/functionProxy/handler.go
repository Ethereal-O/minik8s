package functionProxy

import (
	"fmt"
	"github.com/labstack/echo"
	"minik8s/pkg/client"
	"minik8s/pkg/util/config"
	"net/http"
	"strings"
)

func forwardRequest(c2 echo.Context) error {
	formParams, _ := c2.FormParams()
	formData := make(map[string]string)
	for key, values := range formParams {
		value := strings.Join(values, ",")
		formData[key] = value
	}

	funcName := formData["function"]
	tarFunction := client.GetFunction(funcName)
	if tarFunction == nil {
		return c2.String(http.StatusOK, "The function not exist!")
	} else if tarFunction.Runtime.Status != config.RUNNING_STATUS {
		activate(config.FUNC_NAME + "_rs_" + tarFunction.FaasName)
		return c2.String(http.StatusOK, "The function not up!")
	} else {
		addFlow(config.FUNC_NAME + "_service_" + tarFunction.FaasName)
		formData["module"] = tarFunction.Module
		url := "http://" + tarFunction.Runtime.FunctionIp + ":8081" + "/run"
		res := client.ForwardPostData(formData, url)
		fmt.Println(res)
		return c2.String(http.StatusOK, res)
	}
}
