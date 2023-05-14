package request

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"minik8s/pkg/exeFile"
	"net/http"
	"net/url"
	"strings"
)

var file string

var RequestCmd = &cobra.Command{
	Use:   "request",
	Short: "request for the function of serverless",
	Long:  "this is the main cmd to request for the function of serverless",
	Run:   doit,
}

func doit(cmd *cobra.Command, args []string) {
	params := exeFile.ReadRequest(file)
	urlMap := url.Values{}
	for key, value := range params {
		urlMap.Add(key, value)
	}

	request, err := http.NewRequest("POST", params["url"], strings.NewReader(urlMap.Encode()))
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
	fmt.Println(string(b))
}

func init() {
	RequestCmd.Flags().StringVarP(&file, "file", "f", "", "Path to yaml file")
	RequestCmd.MarkFlagRequired("file")
}

func Request() *cobra.Command {
	return RequestCmd
}
