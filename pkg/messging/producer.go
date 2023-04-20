package messging

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	ec "go.etcd.io/etcd/client/v3"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/stringParse"
)

func Producer(key string, c chan *ec.Event) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	nsqAddress := config.NSQ_PEODUCER
	topic := stringParse.Reform(key)
	go func() {
		config := nsq.NewConfig()
		producer, err := nsq.NewProducer(nsqAddress, config)
		if err != nil {
			fmt.Printf("[ERROR] create producer failed, err:%v\n", err)
		}
		fmt.Printf("[nsq] producer start!(%s)\n", topic)
		var s string
		for {
			select {
			case ev := <-c:
				s = fmt.Sprintf("Type: %s\nKey:%s\nValue:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				fmt.Println("[nsq]   ", s)
				err := producer.Publish(stringParse.Reform(key), []byte(ev.Kv.Value))
				if err != nil {
					fmt.Println("[nsq]" + err.Error())
				}
			case <-ctx.Done():
				producer.Stop()
				fmt.Printf("[nsq] producer end!(%s)\n", topic)
				return
			}
		}
	}()
	return cancel
}
