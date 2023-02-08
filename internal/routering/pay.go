package routering

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/pay"
)

type alipay struct {
}

func (alipay) Init(ctx *gin.RouterGroup) {
	PayRouter := ctx.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}
