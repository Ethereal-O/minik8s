package etcd

import (
	"context"
	"fmt"
	ec "go.etcd.io/etcd/client/v3"
	"minik8s/pkg/util/config"
	"os"
	"time"
)

func Get_etcdClient() *ec.Client {
	cli, err := ec.New(ec.Config{
		Endpoints:   []string{config.ETCD_ENDPOINT},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	return cli
}

func Set_etcd(key string, value string) error {
	cli := Get_etcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	_, err := cli.Put(ctx, key, value)

	defer cli.Close()
	defer cancel()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}
	return nil
}

func Del_etcd(key string) error {
	cli := Get_etcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	_, err := cli.Delete(ctx, key)

	defer cli.Close()
	defer cancel()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}
	return nil
}

func Get_etcd(key string, withPrix bool) []string {
	cli := Get_etcdClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	kv := ec.NewKV(cli)
	var resp *ec.GetResponse
	var err error
	if withPrix {
		resp, err = kv.Get(ctx, key, ec.WithPrefix())
	} else {
		resp, err = kv.Get(ctx, key)
	}

	defer cli.Close()
	defer cancel()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	res := make([]string, 0, 0)
	for _, item := range resp.Kvs {
		res = append(res, string(item.Value))
	}
	return res
}

func Watch_etcd(key string, prix bool, c chan *ec.Event) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cli := Get_etcdClient()
		var rch ec.WatchChan
		if prix {
			rch = cli.Watch(context.Background(), key, ec.WithPrefix())
		} else {
			rch = cli.Watch(context.Background(), key)
		}
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("[etcd]   Type: %s\nKey:%s\nValue:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				c <- ev
			}
		}
		<-ctx.Done()
		cli.Close()
	}()
	return cancel
}
