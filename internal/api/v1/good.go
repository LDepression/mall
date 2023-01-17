package v1

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"strconv"
)

type good struct {
}

func NewGood() *good {
	return &good{}
}
func (*good) SearchGoodsList(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var req form.GoodsFilterRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	goodsInfo, err1 := logic.Group.Good.GetGoodsList(req)
	if err1 != nil {
		rly.Reply(err1)
		return
	}
	rly.Reply(nil, goodsInfo)
}

func (*good) CreateGood(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var req form.CreateGoodReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	goodsInfo, err1 := logic.Group.Good.CreateGood(ctx, req)
	if err1 != nil {
		rly.Reply(err1)
		return
	}
	rly.Reply(nil, goodsInfo)
}

func (*good) DeleteGood(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	err1 := logic.Group.Good.DeleteGood(ctx, int32(id))
	if err1 != nil {
		rly.Reply(err1)
		return
	}
	rly.Reply(nil)
}

func (*good) UpdateGood(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var updateReply form.UpdateGood
	if err := ctx.ShouldBind(&updateReply); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	goodInfo, err := logic.Group.Good.UpdateGood(ctx, updateReply)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, goodInfo)
}

func (*good) GetGoodByID(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var ReqID form.GetGoodID
	if err := ctx.ShouldBindQuery(&ReqID); err != nil {
		base.HandleValidatorError(ctx, err)
	}
	replyInfo, err := logic.Group.Good.GetGoodByID(ctx, ReqID.ID)
	if err != nil {
		rly.Reply(err)
	}
	rly.Reply(nil, replyInfo)
}
