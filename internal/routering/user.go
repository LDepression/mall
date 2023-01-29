package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
	"mall/internal/middleware"
)

type user struct{}

func (user) Init(ctx *gin.RouterGroup) {

	user := v1.NewUser()
	ctx.POST("/login", user.Login)
	ctx.POST("/register", user.Register)
	group := ctx.Group("/u")
	auth1 := group.Use(middleware.Auth())
	auth2 := group.Use(middleware.ManagerAuth(), middleware.Auth())
	{
		auth1.GET("/:id", user.GetUserByID)
		auth1.POST("/userInfo", user.UpdateUserInfo)
		auth1.POST("/updateEmail", user.UpdateEmail)
		auth1.POST("/modifyPassword", user.ModifyPassword)
		auth2.POST("/refresh", user.RefreshToken)
		auth2.GET("/list", user.GetUsers)
		auth2.DELETE("/delete", user.Delete)
		auth1.POST("/searchUserByName", user.SearchUsers)
	}
}
