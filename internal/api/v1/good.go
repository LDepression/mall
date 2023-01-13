package v1

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
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

//func (*good) CreateGood(ctx *gin.Context) {
//	rly := app.NewResponse(ctx)
//	var req form.CreateGoodReq
//	err := ctx.ShouldBind(&req)
//	if err != nil {
//		base.HandleValidatorError(ctx, err)
//		return
//	}
//	goodsInfo, err1 := logic.Group.Good.CreateGood(req)
//	if err1 != nil {
//		rly.Reply(err1)
//		return
//	}
//	rly.Reply(nil, goodsInfo)
//}
