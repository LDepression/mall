package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
)

type email struct {
}

func (email) Init(ctx *gin.RouterGroup) {
	ctx.POST("/sendEmail", v1.SendEmailCode)
}
