package routering

import (
	v1 "mall/internal/api/v1"
	"mall/internal/global"
	"mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

type email struct {
}

func (email) Init(ctx *gin.RouterGroup) {
	email := ctx.Group("email", middleware.LimitAPI(GetLimiters(global.Setting.Limit.APILimit.Email)))
	email.POST("/sendEmail", v1.SendEmailCode)
}
