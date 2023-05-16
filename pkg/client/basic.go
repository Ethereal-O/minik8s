package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/pkg/util/config"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ForwardPostData(formData map[string]string, tagetUrl string) string {
	urlMap := url.Values{}
	for key, value := range formData {
		urlMap.Add(key, value)
	}

	request, err := http.NewRequest("POST", tagetUrl, strings.NewReader(urlMap.Encode()))
	if err != nil {
		panic(err)
	}
	fmt.Println("request.url: ", request.URL)
	fmt.Println("request.method: ", request.Method)

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//发送请求给服务端
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//服务端返回数据
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func postFormData(key string, prix bool, crt string) string {
	urlMap := url.Values{}
	urlMap.Add("key", key)
	if prix {
		urlMap.Add("prix", "true")
	} else {
		urlMap.Add("prix", "false")
	}
	urlMap.Add("crt", crt)

	request, err := http.NewRequest("POST", config.APISERVER_URL, strings.NewReader(urlMap.Encode()))
	if err != nil {
		panic(err)
	}
	fmt.Println("request.url: ", request.URL)
	fmt.Println("request.method: ", request.Method)

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//发送请求给服务端
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//服务端返回数据
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func delete(url string) string {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
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
	return string(b)

}

func get(url string) []string {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var ans []string
	err = json.NewDecoder(res.Body).Decode(&ans)
	if err != nil {
		panic(err)
	}
	return ans
}

func put(url string, dataStr string) string {
	request, err := http.NewRequest("PUT", url, strings.NewReader(dataStr))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
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
	return string(b)
}
