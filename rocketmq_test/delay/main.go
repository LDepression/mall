package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.28.11:9876"}))
	if err != nil {
		panic("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	msg := primitive.NewMessage("mall", []byte("this is delay msg"))
	msg.WithDelayTimeLevel(3)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Println("发送失败", err)
	} else {
		fmt.Println("发送成功,res:", res.String())
	}
	if err := p.Shutdown(); err != nil {
		panic(err)
	}
	//支付的时候 淘宝,购票,12306的时候一定要有超时归还的逻辑
	/*
		但是为什么不能去轮询呢?
		比如说我的定时任务是在12:00执行了一次了,用户在12:01下了单.应该是12:31超时.但是定时任务会在12:30执行一次,发现没有超时.
		再一次执行要是13:00了,此时我12:01的订单只能是13:00将其设置为超时.那么多了29分钟.这就很严重了

		但是要是我把定时任务的时间设置为1分钟一次的话,那么可能很多次轮询都是无用功

		所以我们要用rocketmq的延迟消息
	*/
}
