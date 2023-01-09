package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mall/internal/dao/db/query"
	"mall/internal/global"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rly := app.NewResponse(ctx)
		token := ctx.GetHeader(global.Setting.Token.AuthType)
		payLoad, err := global.Maker.VerifyToken(token)
		if err != nil {
			rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
			ctx.Abort()
		}
		quser := query.NewUser()
		userInfo, err := quser.GetUserByID(int32(payLoad.UserID))
		if err != nil {
			rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
			ctx.Abort()
		}
		if userInfo.Role == 0 {
			rly.Reply(errcode.ErrServer.WithDetails(errors.New("请先登陆").Error()))
			ctx.Abort()
		}
		ctx.Set(global.Setting.Token.AuthKey, gin.H{
			"payload": payLoad,
			"role":    userInfo.Role,
		})
		ctx.Next()
	}
}
