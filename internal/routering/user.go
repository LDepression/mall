package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
)

type user struct{}

func (user) Init(group *gin.RouterGroup) {

	user := v1.NewUser()
	group.POST("/login", user.Login)
	group.POST("/register", user.Register)
}
