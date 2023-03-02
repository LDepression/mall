/**
 * @Author: lenovo
 * @Description:
 * @File:  limit
 * @Version: 1.0.0
 * @Date: 2023/02/26 14:56
 */

package middleware

import (
	"context"
	"errors"
	"mall/internal/dao"
	"mall/internal/dao/redis/query"
	"mall/internal/global"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
	limit "mall/internal/pkg/limiter/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//LimitIP 对IP进行限流
func LimitIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rly := app.NewResponse(ctx)
		IP := ctx.RemoteIP()
		bucket, err := dao.Group.Redis.GetBucket(ctx, query.BucketRequest{
			Key:     IP,
			Cap:     global.Setting.Limit.IPLimit.Cap,
			GenNum:  global.Setting.Limit.IPLimit.GenNum,
			GenTime: global.Setting.Limit.IPLimit.GenTime,
			Cost:    global.Setting.Limit.IPLimit.Cost,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				rly.Reply(errcode.ErrTimeOut)
			} else {
				rly.Reply(errcode.ErrServer, "内部错误")
				zap.S().Info(err.Error())
			}
			ctx.Abort()
			return
		}
		if bucket.Success != true {
			rly.Reply(errcode.ErrSendTooMany)
			zap.S().Info("ip limit :", zap.String("ip", IP))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

//LimitAPI 对于某一个API进行限流
func LimitAPI(limit limit.RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rly := app.NewResponse(ctx)
		if err := limit.Wait(ctx); err != nil {
			zap.S().Info("api limit error: %v", err)
			rly.Reply(errcode.ErrTimeOut)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
