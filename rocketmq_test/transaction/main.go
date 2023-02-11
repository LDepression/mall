/**
 * @Author: lenovo
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2023/02/09 16:00
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type orderListener struct {
}

func (l *orderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 5)
	fmt.Println("本地逻辑执行失败")
	/*
		本地执行逻辑没有问题
	*/
	return primitive.UnknowState
}

//回查机制
func (l *orderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("开始执行本地回查")
	time.Sleep(time.Second * 5)
	/*
		本地执行逻辑没有问题
	*/
	return primitive.CommitMessageState
}
func main() {
	p, err := rocketmq.NewTransactionProducer(
		&orderListener{},
		producer.WithNameServer([]string{"192.168.28.13:9876"}),
	)
	if err != nil {
		panic("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	if err := p.Start(); err != nil {
		panic(err)
	}
	result, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TopicTrans", []byte("this is transaction3")))
	if err != nil {
		fmt.Println("发送消息失败 err:", err)
	} else {
		fmt.Println("发送消息成功:", result.String())
	}
	time.Sleep(time.Hour)
	if err := p.Shutdown(); err != nil {
		panic(err)
	}
}
