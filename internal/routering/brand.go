package routering

import (
	"github.com/gin-gonic/gin"
	v1 "mall/internal/api/v1"
	"mall/internal/logic"
	"sync"
)

type brand struct {
}

func (brand) Init(ctx *gin.RouterGroup) {
	group := ctx.Group("/b")
	logic.Group.Brand.Lock = new(sync.Mutex)
	brand := v1.NewBrand()
	group.GET("/:id", brand.SearchBrand)
	group.POST("", brand.CreateBrand)
	group.PUT("/:id", brand.UpdateBrand)
	group.DELETE("/:id", brand.DeleteBrand)
	group.GET("", brand.BrandList)
}
