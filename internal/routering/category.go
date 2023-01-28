package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
)

type category struct {
}

func (category) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/c")
	category := v1.NewCategory()
	group.GET("", category.GetAllCategoryList)
	group.GET("/:id", category.SearchCategory)
	group.POST("", category.CreateCategory)
	group.PUT("/:id", category.UpdateCategory)
	group.DELETE("/:id", category.DeleteCategory)
}
