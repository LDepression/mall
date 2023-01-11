package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
)

type upload struct {
}

func (upload) Init(ctx *gin.RouterGroup) {
	ctx.POST("/upload", v1.UploadFile)
}
