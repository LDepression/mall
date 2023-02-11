package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.28.13:9876"}))
	if err != nil {
		panic("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	res, err := p.SendSync(context.Background(), primitive.NewMessage("mall", []byte("this is mall")))
	if err != nil {
		fmt.Println("发送失败", err)
	} else {
		fmt.Println("发送成功,res:", res.String())
	}
	if err := p.Shutdown(); err != nil {
		panic(err)
	}
}
