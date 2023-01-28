package v1

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"strconv"
)

type brand struct {
}

func NewBrand() *brand {
	return &brand{}
}

func (*brand) CreateBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var createBrand form.CreateBrand
	if err := ctx.ShouldBindJSON(&createBrand); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.Brand.CreateBrand(createBrand); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "增加品牌成功")
}
func (*brand) DeleteBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err := logic.Group.Brand.DeleteBrand(i); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "删除品牌成功")
}
func (*brand) UpdateBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var updateBrand form.UpdateBrand
	if err := ctx.ShouldBindJSON(&updateBrand); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err := logic.Group.Brand.UpdateBrand(i, updateBrand); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "更新商品成功")
}
func (*brand) SearchBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	BrandInfo, err := logic.Group.Brand.GetBrandByID(i)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, BrandInfo)
}

func (*brand) BrandList(ctx *gin.Context) {
	var reqBrandList form.ReqBrandsList
	if err := ctx.ShouldBindQuery(&reqBrandList); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	rly := app.NewResponse(ctx)
	BrandListInfo, err := logic.Group.Brand.BrandList(reqBrandList)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.ReplyList(nil, BrandListInfo)
}
