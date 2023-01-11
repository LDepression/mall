package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/internal/dao/db/query"
	"mall/internal/global"
	"mall/internal/model/common"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rly := app.NewResponse(ctx)
		Token := ctx.GetHeader(global.Setting.Token.AuthType)
		payLoad, err := global.Maker.VerifyToken(Token)
		if err != nil {
			rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
			ctx.Abort()
			return
		}
		quser := query.NewUser()
		userInfo, err := quser.GetUserByID(payLoad.UserID)
		if err != nil {
			rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
			ctx.Abort()
			return
		}
		if userInfo.Role == 0 {
			rly.Reply(errcode.ErrNotLogin)
			ctx.Abort()
			return
		}
		ctx.Set(global.Setting.Token.AuthKey, common.PalLoad{
			Role:    userInfo.Role,
			PalLoad: *payLoad,
		})
		ctx.Next()
	}
}

func GetPayload(ctx *gin.Context) (common.PalLoad, errcode.Err) {
	payload, ok := ctx.Get(global.Setting.Token.AuthKey)
	if !ok {
		return payload.(common.PalLoad), errcode.ErrNotLogin
	}
	return payload.(common.PalLoad), nil
}

func ManagerAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rly := app.NewResponse(ctx)
		payload, err := GetPayload(ctx)
		if err != nil {
			zap.S().Info("用户未登录")
			rly.Reply(err)
			ctx.Abort()
			return
		}
		if payload.Role == 0 || payload.Role == 1 {
			zap.S().Info("用户权限不足")
			rly.Reply(errcode.ErrNotManager)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
