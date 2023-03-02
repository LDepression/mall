/**
 * @Author: lenovo
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2023/02/09 15:29
 */

package main

import (
	"fmt"
	"mall/internal/logic"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.28.16:9876"}),
		consumer.WithGroupName("mall-reback11"),
	)
	if err := c.Subscribe("order_back", consumer.MessageSelector{}, logic.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	time.Sleep(time.Hour) //不能让主携程退出
	_ = c.Shutdown()
}
