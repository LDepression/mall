package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
	"strconv"
)

type Category struct {
}

func NewCategory() *Category {
	return &Category{}
}

func (*Category) CreateCategory(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var categoryForm form.CreateCategory
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err1 := logic.Group.Category.CreateCategory(categoryForm); err1 != nil {
		rly.Reply(err1)
		return
	}
	rly.Reply(nil, "增加分类成功")
}

func (*Category) DeleteCategory(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	i, err := strconv.Atoi(idStr)
	if err != nil {
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	if err1 := logic.Group.Category.DeleteCategory(int32(i)); err1 != nil {
		rly.Reply(err1)
		return

	}
	rly.Reply(nil, "删除分类成功")
}

func (*Category) UpdateCategory(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var categoryForm form.UpdateCategory
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err1 := logic.Group.Category.UpdateCategory(categoryForm, i); err1 != nil {
		rly.Reply(err1)
		return

	}
	rly.Reply(nil, "更新分类成功")
}

func (*Category) SearchCategory(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	replyInfo, err1 := logic.Group.Category.SearchCategory(i)
	if err1 != nil {
		rly.Reply(err1)
		return

	}
	rly.Reply(nil, replyInfo)
}
func (*Category) GetAllCategoryList(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	replyInfo, err1 := logic.Group.Category.GetAllCategoryList()
	if err1 != nil {
		rly.Reply(err1)
		return

	}
	data := make([]interface{}, 0)
	err := json.Unmarshal([]byte(replyInfo.JsonData), &data)
	if err != nil {
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	rly.Reply(nil, data)
}
