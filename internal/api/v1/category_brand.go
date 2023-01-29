package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
	"strconv"
)

type CategoryBrand struct {
}

func NewCategoryBrand() *CategoryBrand {
	return &CategoryBrand{}
}

func (c *CategoryBrand) GetCategoryBrandList(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	listReq := form.CategoryBrandList{}
	if err := ctx.ShouldBindQuery(&listReq); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := logic.Group.CategoryBrand.CategoryBrandList(listReq.Page, listReq.PagePerNum)
	if err != nil {
		rly.Reply(errcode.ErrServer.WithDetails("查询失败"))
		return
	}
	//此时创建一个map就好了,练习一下自定义数据格式给前端
	reMap := make(map[string]interface{})
	reMap["total"] = rsp.Total
	result := make([]interface{}, 0)
	for _, categoryBrandInfo := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = categoryBrandInfo.ID
		reMap["category"] = map[string]interface{}{
			"category_id": categoryBrandInfo.Category.ID,
			"name":        categoryBrandInfo.Category.Name,
			"level":       categoryBrandInfo.Category.Level,
		}
		reMap["brand"] = map[string]interface{}{
			"brand_id": categoryBrandInfo.Brand.ID,
			"name":     categoryBrandInfo.Brand.Name,
			"logo":     categoryBrandInfo.Brand.Logo,
		}
		result = append(result, reMap)
	}
	reMap["data"] = result
	rly.Reply(nil, reMap)
}
func (c *CategoryBrand) NewCategoryBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var CategoryBrandRequest form.CreateCategoryBrand
	if err := ctx.ShouldBindJSON(&CategoryBrandRequest); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := logic.Group.CategoryBrand.CreateCategoryBrand(CategoryBrandRequest)
	if err != nil {
		zap.S().Info("logic.Group.CategoryBrand failed", err)
		rly.Reply(errcode.ErrServer.WithDetails("你输入的关联字段不存在"))
		return
	}
	re := make(map[string]interface{})
	re["id"] = rsp.ID
	rly.Reply(nil, re)
}
func (c *CategoryBrand) UpdateCategoryBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var updateRequest form.UpdateCategoryBrand
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err := logic.Group.CategoryBrand.UpdateCategoryBrand(i, updateRequest); err != nil {
		zap.S().Info("logic.Group.CategoryBrand.UpdateCategoryBrand failed err:", err)
		rly.Reply(errcode.ErrServer.WithDetails("你输入的关联字段不存在"))
		return
	}
	rly.Reply(nil, "更新成功")
}
func (c *CategoryBrand) DeleteCategoryBrand(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err := logic.Group.CategoryBrand.DeleteCategoryBrand(i); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "删除成功")
}
