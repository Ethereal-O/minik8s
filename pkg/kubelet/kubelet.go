package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/messging"
	"minik8s/pkg/object"
	"minik8s/pkg/util/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start_kubelet() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)

	podChan, podStop := messging.Watch("/"+config.POD_TYPE, true)
	go dealPod(podChan)

	<-c
	podStop()
	time.Sleep(2 * time.Second)
	return
}

func dealPod(podChan chan string) {
	for {
		select {
		case mes := <-podChan:
			//fmt.Println("[this]", mes)
			var tarPod object.Pod
			err := json.Unmarshal([]byte(mes), &tarPod)
			if err != nil {
				fmt.Println(err.Error())
			}
			CreatePod(&tarPod)
		}
	}
}
