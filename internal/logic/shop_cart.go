package logic

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
)

type shopCart struct {
}

func (c *shopCart) CartItemList(userId int64) (*reply.CartItemListResponse, errcode.Err) {
	rsp := &reply.CartItemListResponse{}
	var count int64
	dao.Group.DB.Model(&model.ShoppingCart{}).Where("user=?", userId).Count(&count)
	if count == 0 {
		return nil, errcode.ErrNotFound
	}
	rsp.Total = count
	cartItems, err := query.Group.ShopCart.GetShopCartByUserID(userId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrServer
	}

	for _, cartItem := range cartItems {
		rsp.Data = append(rsp.Data, &reply.ShopCartInfoResponse{
			ID:      cartItem.ID,
			UserID:  cartItem.User,
			GoodID:  cartItem.Goods,
			Num:     cartItem.Nums,
			Checked: cartItem.Checked,
		})
	}
	return rsp, nil
}
func (c *shopCart) CreateCartItem(req form.ShopCartItemForm, userID int64) errcode.Err {
	var shopCart model.ShoppingCart
	if result := dao.Group.DB.Where(model.ShoppingCart{User: int32(userID), Goods: req.GoodsId}).Find(&shopCart); result.RowsAffected == 0 {
		if err := query.Group.ShopCart.CreateShopCart(model.ShoppingCart{
			User:    int32(userID),
			Goods:   req.GoodsId,
			Nums:    req.Nums,
			Checked: false,
		}); err != nil {
			zap.S().Info("query.Group.ShopCart.CreateShopCart,err:", err)
			return errcode.ErrServer
		}
	} else {
		shopCart.Nums += req.Nums
		if err := query.Group.ShopCart.UpdateShopCart(shopCart); err != nil {
			zap.S().Info("query.Group.ShopCart.UpdateShopCart,err:", err)
			return errcode.ErrServer
		}
	}
	return nil
}
func (c *shopCart) DeleteCartItem(userID int32, GoodID int) errcode.Err {
	if err := query.Group.ShopCart.DeleteShopCart(userID, int32(GoodID)); err != nil {
		zap.S().Info("DeleteShopCart failed", err)
		return errcode.ErrServer
	}
	return nil
}
func (c *shopCart) UpdateCartItem(userID, goodID int32, updateForm form.ShopCartItemUpdateForm) errcode.Err {
	var req model.ShoppingCart
	if updateForm.Checked != nil {
		req.Checked = *updateForm.Checked
	}
	req.User = userID
	req.Goods = goodID
	if updateForm.Nums > 0 {
		req.Nums = updateForm.Nums
	}
	if err := query.Group.ShopCart.UpdateShopCart(req); err != nil {
		zap.S().Info("query.Group.ShopCart.UpdateShopCart,err:", err)
		return errcode.ErrServer
	}
	return nil
}
