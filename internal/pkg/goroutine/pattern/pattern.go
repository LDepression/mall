package pattern

import (
	"context"
)

// Or 监听多个channel 只要有一个返回消息就返回
func Or(ctxs ...context.Context) context.Context {
	switch len(ctxs) {
	case 0:
		return nil
	case 1:
		return ctxs[0]
	}
	orCtx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		switch len(ctxs) {
		case 2:
			select {
			case <-ctxs[0].Done():
			case <-ctxs[1].Done():
			}
		default:
			select {
			case <-ctxs[0].Done():
			case <-ctxs[1].Done():
			case <-ctxs[2].Done():
			case <-Or(append(ctxs[3:], orCtx)...).Done(): // 递归退出
			}
		}
	}()
	return orCtx
}
