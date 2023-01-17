package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
	"mall/internal/logic"
	"sync"
)

type good struct {
}

func (good) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/g")
	good := v1.NewGood()
	logic.Group.Good.Lock = new(sync.Mutex)
	group.POST("/list", good.SearchGoodsList)
	group.POST("", good.CreateGood)
	group.DELETE("/:id", good.DeleteGood)
	group.PUT("", good.UpdateGood)
	group.GET("", good.GetGoodByID)
}
