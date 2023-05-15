package request

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
	"minik8s/pkg/util/config"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var paramFile string
var funcName string

var RequestCmd = &cobra.Command{
	Use:   "request",
	Short: "request for the function of serverless",
	Long:  "this is the main cmd to request for the function of serverless",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	tarFunction := client.GetFunction(funcName)
	if tarFunction == nil {
		fmt.Println("The function not exist!")
		return
	}
	if tarFunction.Runtime.Status != config.RUNNING_STATUS {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println("Waiting for the pod to start (cold start)!")
			tarFunction = client.GetFunction(funcName)
			if tarFunction.Runtime.Status == config.RUNNING_STATUS {
				break
			}
		}
	}
	params := exeFile.ReadRequest(paramFile)
	urlMap := url.Values{}
	for key, value := range params {
		urlMap.Add(key, value)
	}
	urlMap.Add("function", tarFunction.FuncName)
	urlMap.Add("module", tarFunction.Module)
	url := "http://" + tarFunction.Runtime.FunctionIp + ":8081" + "/run"
	request, err := http.NewRequest("POST", url, strings.NewReader(urlMap.Encode()))
	if err != nil {
		panic(err)
	}
	fmt.Println("request.url: ", request.URL)
	fmt.Println("request.method: ", request.Method)
	fmt.Println("params: ", params)

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("The result is:", string(b))
}

func init() {
	RequestCmd.Flags().StringVarP(&funcName, "funcName", "f", "", "The function name")
	RequestCmd.Flags().StringVarP(&paramFile, "paramFile", "p", "", "Path to the param yaml file")
	RequestCmd.MarkFlagRequired("funcName")
	RequestCmd.MarkFlagRequired("paramFile")
}

func Request() *cobra.Command {
	return RequestCmd
}
