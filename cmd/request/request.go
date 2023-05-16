package request

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/client"
	"minik8s/pkg/exeFile"
	"minik8s/pkg/util/config"
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
	params := exeFile.ReadRequest(paramFile)
	params["function"] = funcName
	url := config.FUNCTION_PROXY_URL
	res := client.ForwardPostData(params, url)
	if res == "The function not exist!" {
		fmt.Println("The function not exist!")
	} else if res == "The function not up!" {
		for {
			time.Sleep(3 * time.Second)
			fmt.Println("Waiting for the pod to start (cold start)!")
			res = client.ForwardPostData(params, url)
			if res != "The function not up!" {
				time.Sleep(5 * time.Second)
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
	RequestCmd.Flags().StringVarP(&paramFile, "paramFile", "p", "", "Path to the param yaml file")
	RequestCmd.MarkFlagRequired("funcName")
	RequestCmd.MarkFlagRequired("paramFile")
}

func Request() *cobra.Command {
	return RequestCmd
}
