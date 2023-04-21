package messging

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"minik8s/pkg/util/config"
	"minik8s/pkg/util/stringParse"
	"time"
)

type MyHandler struct {
	myChan chan string
}

func (m *MyHandler) HandleMessage(msg *nsq.Message) (err error) {
	m.myChan <- string(msg.Body)
	return
}

func initConsumer(ctx context.Context, topic string, channel string, address string, resChan chan string) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 30 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return
	}
	consumer := &MyHandler{
		myChan: resChan,
	}
	c.AddHandler(consumer)

	if err := c.ConnectToNSQLookupd(address); err != nil {
		fmt.Println("[nsq]" + err.Error())
	} else {
		fmt.Printf("[nsq] consumer start!(%s)\n", topic)
	}
	<-ctx.Done()
	c.Stop()
	fmt.Printf("[nsq] consumer end!(%s)\n", topic)
	return
}

func Consumer(key string, crt string, resChan chan string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go initConsumer(ctx, stringParse.Reform(key), crt, config.NSQ_CONSUMER, resChan)
	return cancel
}
