package routering

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
)

type abase struct {
}

func (abase) Init(ctx *gin.RouterGroup) {
	ctx.GET("captcha", base.GetCaptcha)
}
