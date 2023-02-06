package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
	"mall/internal/middleware"
)

type shopCart struct {
}

func (shopCart) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/shopCart").Use(middleware.Auth())
	shopCart := v1.NewShopCart()
	group.POST("", shopCart.New)
	group.DELETE("/:id", shopCart.Delete)
	group.PUT("/:id", shopCart.Update)
	group.GET("", shopCart.List)
}
