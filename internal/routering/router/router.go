package router

import (
	"mall/internal/middleware"
	"mall/internal/routering"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery(true), middleware.LogBody())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	root := r.Group("api/v1", middleware.LimitIP())
	routering.Group.User.Init(root)
	routering.Group.Base.Init(root)
	routering.Group.Email.Init(root)
	routering.Group.Upload.Init(root)
	routering.Group.Good.Init(root)
	routering.Group.Category.Init(root)
	routering.Group.Brand.Init(root)
	routering.Group.ShopCart.Init(root)
	routering.Group.Order.Init(root)
	routering.Group.AliPay.Init(root)
	return r
}
