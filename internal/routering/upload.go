package routering

import (
	v1 "mall/internal/api/v1"
	"mall/internal/global"
	"mall/internal/middleware"

	"github.com/gin-gonic/gin"
)

type upload struct {
}

func (upload) Init(ctx *gin.RouterGroup) {
	//GetLimiters GetAPI 返回的是一个混合的限流器,并且是进过了排序之后的限流器
	upload := ctx.Group("upload", middleware.LimitAPI(GetLimiters(global.Setting.Limit.APILimit.Upload)))
	upload.POST("/upload", v1.UploadFile)
}
