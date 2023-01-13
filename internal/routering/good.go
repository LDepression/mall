package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
)

type good struct {
}

func (good) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/g")
	good := v1.NewGood()
	group.POST("/list", good.SearchGoodsList)
	//group.POST("", good.CreateGood)
	//group.DELETE("/:id", good.DeleteGood)
	//group.PUT("/:id", good.UpdateGood)
}
