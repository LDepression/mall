package heal

import (
	"context"
	"fmt"
	"log"
	"time"

	"ttms/internal/pkg/goroutine/pattern"
)

/*
	管理协程并且在协程异常时重新启动协程
*/

// StartGoroutineFn
// 创建一个可以监控和重启的goroutine的方式
// 参数：退出channel,心跳时间
// 返回值：返回心跳的channel
type StartGoroutineFn func(ctx context.Context, pulseInterval time.Duration) <-chan struct{}

// NewSteward
// 新建一个管理员
// 参数：下游的超时时间，创建一个可以监控和重启的goroutine的方式
// 返回值：返回一个创建一个受管理的goroutine和其管理者的函数的创建方式
func NewSteward(name string, timeout time.Duration, startGoroutine StartGoroutineFn) StartGoroutineFn {
	return func(ctx context.Context, pulseInterval time.Duration) <-chan struct{} {
		heartBeat := make(chan struct{})
		go func() {
			defer close(heartBeat)
			var (
				// 管理者用于通知下游退出的channel
				wardCtx context.Context
				cancel  context.CancelFunc
			)
			var wardHeartbeat <-chan struct{} // 管理员用于接收下游心跳的channel
			startWard := func() {
				log.Println("stewart: start new goroutine:", name)
				wardCtx, cancel = context.WithCancel(ctx)                           // 初始化退出channel
				wardHeartbeat = startGoroutine(pattern.Or(wardCtx, ctx), timeout/2) // 启动下游，其心跳间隔是超时间隔的一半
			}
			startWard()                            // 启动受监管的goroutine
			pulse := time.NewTicker(pulseInterval) // 定时回复上游的心跳
			defer pulse.Stop()
		monitorLoop:
			for {
				timeoutSignal := time.After(timeout) // 用来提醒自己下游超时了
				for {
					select {
					case <-pulse.C: // 回复心跳
						select {
						case heartBeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat: // 接收到下游的心跳则继续监视
						continue monitorLoop
					case t := <-timeoutSignal: // 没收到下游的心跳则重启下游
						log.Println(fmt.Sprintf("timeout signal: name:%s time:%v", name, t))
						cancel()
						startWard() // 使用之前的方式重新启动下游
						continue monitorLoop
					case <-ctx.Done():
						return
					}
				}
			}
		}()
		return heartBeat
	}
}
