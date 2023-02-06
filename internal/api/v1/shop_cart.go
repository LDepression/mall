package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/middleware"
	"mall/internal/pkg/app"
	"strconv"
)

type shopCart struct {
}

func NewShopCart() *shopCart {
	return &shopCart{}
}

func (c *shopCart) New(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	payLoad, ok := middleware.GetPayload(ctx)
	if ok != nil {
		zap.S().Info("GetPayload failed,err:", ok)
		rly.Reply(ok)
		return
	}
	userId := payLoad.PalLoad.UserID
	var req form.ShopCartItemForm
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.ShopCart.CreateCartItem(req, userId); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "添加购物车成功")
}
func (c *shopCart) Delete(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	payLoad, ok := middleware.GetPayload(ctx)
	if ok != nil {
		zap.S().Info("GetPayload failed,err:", ok)
		rly.Reply(ok)
		return
	}
	userId := payLoad.PalLoad.UserID
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	if err := logic.Group.ShopCart.DeleteCartItem(int32(userId), i); err != nil {
		zap.S().Info("删除购物车失败")
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "删除成功")
}
func (c *shopCart) Update(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	payLoad, ok := middleware.GetPayload(ctx)
	if ok != nil {
		zap.S().Info("GetPayload failed,err:", ok)
		rly.Reply(ok)
		return
	}
	userId := payLoad.PalLoad.UserID
	var req form.ShopCartItemUpdateForm
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	idStr := ctx.Param("id")
	goodID, _ := strconv.Atoi(idStr)
	if err := logic.Group.ShopCart.UpdateCartItem(int32(userId), int32(goodID), req); err != nil {
		zap.S().Info("logic.Group.ShopCart.UpdateCartItem failed,err:", err)
		rly.Reply(err)
		return
	}
	rly.Reply(nil, "更新成功")
}
func (c *shopCart) List(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	payLoad, ok := middleware.GetPayload(ctx)
	if ok != nil {
		zap.S().Info("GetPayload failed,err:", ok)
		rly.Reply(ok)
		return
	}
	userId := payLoad.PalLoad.UserID
	rsp, err := logic.Group.ShopCart.CartItemList(userId)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, rsp)
}
