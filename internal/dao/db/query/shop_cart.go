package query

import (
	"mall/internal/dao"
	"mall/internal/model"
)

type shopCart struct {
}

func (c *shopCart) GetShopCartByUserID(userID int64) ([]*model.ShoppingCart, error) {
	var shopCarts []*model.ShoppingCart
	result := dao.Group.DB.Model(&model.ShoppingCart{}).Where("user=?", userID).Find(&shopCarts)
	return shopCarts, result.Error
}

func (c *shopCart) CreateShopCart(cart model.ShoppingCart) error {
	result := dao.Group.DB.Model(&model.ShoppingCart{}).Create(&cart)
	return result.Error
}

func (c *shopCart) UpdateShopCart(cart model.ShoppingCart) error {
	result := dao.Group.DB.Where("goods=? and user=?", cart.Goods, cart.User).Updates(&cart)
	return result.Error
}
func (c *shopCart) DeleteShopCart(userID, goodID int32) error {
	result := dao.Group.DB.Where("user=? and goods =?", userID, goodID).Delete(&model.ShoppingCart{})
	return result.Error
}

func (c *shopCart) GetShopCartCheckedGoodsByUserID(userID int32) ([]*model.ShoppingCart, error) {
	var shopCarts []*model.ShoppingCart
	result := dao.Group.DB.Model(&model.ShoppingCart{}).Where(&model.ShoppingCart{Checked: true, User: userID}).Find(&shopCarts)
	return shopCarts, result.Error
}
