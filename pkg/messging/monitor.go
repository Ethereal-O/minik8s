package messging

import (
	"minik8s/pkg/client"
	"time"
)

func Watch(key string, prix bool) (chan string, func()) {
	resChan := make(chan string, 20)
	newCrt := client.Post(key, prix, "")
	cancel := Consumer(key, resChan)
	return resChan, func() {
		cancel()
		time.Sleep(1 * time.Second)
		client.Post(key, prix, newCrt)
	}
}
