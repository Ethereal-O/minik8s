package request

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
	"minik8s/pkg/util/config"
	"time"
)

var paramFile string
var funcName string
var dagFile string

var RequestCmd = &cobra.Command{
	Use:   "request",
	Short: "request for the function of serverless",
	Long:  "this is the main cmd to request for the function of serverless",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	if funcName == "" && dagFile == "" {
		fmt.Println("Neither a function nor a workflow!")
		return
	} else if funcName != "" {
		fmt.Println("[A simple request]")
		simpleRequest()
	} else {
		fmt.Println("[A dag workflow]")
		workFlow()
	}

}

func simpleRequest() {
	request := make(map[string]string)
	params := exeFile.ReadRequest(paramFile)
	params["function"] = funcName
	request = params

	url := config.FUNCTION_PROXY_URL + "/simpleRequest"

	doRequest(request, url)
}
func workFlow() {
	request := make(map[string]string)
	workflow := exeFile.ReadWorkFlow(dagFile)
	params := exeFile.ReadRequest(paramFile)
	workflow_json, _ := json.Marshal(workflow)
	params_json, _ := json.Marshal(params)
	request["workflow"] = string(workflow_json)
	request["params"] = string(params_json)

	url := config.FUNCTION_PROXY_URL + "/workflow"

	doRequest(request, url)
}

func doRequest(request map[string]string, url string) {
	res := client.ForwardPostData(request, url)
	if res == "Function not exist!" {
		fmt.Println("Function not exist!")
	} else if res == "Function not up!" {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println("Waiting for the pod(s) to start (cold start)!")
			res = client.ForwardPostData(request, url)
			if res != "Function not up!" {
				fmt.Println("The result is:", res)
				break
			}
		}
	} else {
		fmt.Println("The result is:", res)
	}
}

func init() {
	RequestCmd.Flags().StringVarP(&funcName, "funcName", "f", "", "The function name")
	RequestCmd.Flags().StringVarP(&dagFile, "dagFile", "d", "", "Path to the dag yaml file")
	RequestCmd.Flags().StringVarP(&paramFile, "paramFile", "p", "", "Path to the param yaml file")
	RequestCmd.MarkFlagRequired("paramFile")
}

func Request() *cobra.Command {
	return RequestCmd
}
