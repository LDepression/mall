package router

import (
	"github.com/gin-gonic/gin"
	"mall/internal/middleware"
	"mall/internal/routering"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery(true), middleware.LogBody())
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	root := r.Group("api/v1")
	routering.Group.User.Init(root)
	routering.Group.Base.Init(root)
	routering.Group.Email.Init(root)
	return r
}
