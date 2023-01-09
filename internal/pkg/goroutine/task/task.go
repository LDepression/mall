package task

import (
	"context"
	"log"
	"time"

	"mall/internal/pkg/goroutine/heal"
)

type Task struct {
	Name            string          // task 名
	Ctx             context.Context // 上游ctx
	TaskDuration    time.Duration   // 任务执行周期
	TimeoutDuration time.Duration   // 超时时长
	F               DoFunc          // 执行程序
}
type DoFunc func(parentCtx context.Context)

// NewTickerTask 新建一个被管理的定时任务，并返回可以监听的管理者的心跳
func NewTickerTask(task Task) <-chan struct{} {
	startFun := func(ctx context.Context, pulseInterval time.Duration) <-chan struct{} {
		ticker := time.NewTicker(task.TaskDuration) // 定时任务
		pulse := time.NewTicker(pulseInterval)      // 定期心跳
		heartBeat := make(chan struct{})
		go func() {
			defer ticker.Stop() // 关闭后停止ticker
			defer pulse.Stop()  // 关闭后停止回复心跳
			now := time.Now()
			task.F(ctx) // 起码会执行一次
			log.Println("first exec task:", task.Name, " cost time:", time.Since(now))
			for {
				select {
				case <-ticker.C:
					now := time.Now()
					task.F(ctx)
					log.Println("task: try to exec task:", task.Name, " cost time:", time.Since(now))
				case <-pulse.C:
					heartBeat <- struct{}{} // 必须回复
				case <-ctx.Done():
					log.Println("task: over by stewart:", task.Name)
					return
				}
			}
		}()
		return heartBeat
	}
	return heal.NewSteward(task.Name, task.TimeoutDuration, startFun)(task.Ctx, task.TimeoutDuration)
}
