package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
	"mall/internal/middleware"
)

type order struct {
}

func (order) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/order").Use(middleware.Auth())
	order := v1.NewOrder()
	group.GET("", order.List)
	group.POST("", order.Create)
	group.GET("/:id", order.Details)
	group.PUT("", order.Update)
}
