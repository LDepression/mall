package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/form"
	"mall/internal/model"
)

type inventory struct {
	*gorm.DB
}

func NewInventory() *inventory {
	return &inventory{
		dao.Group.DB,
	}
}

func (i *inventory) CheckGoodID(GoodID int32) error {
	var inv model.Inventory
	result := dao.Group.DB.Where("goods=?", GoodID).Find(&inv)
	return result.Error
}
func (i *inventory) SetInv(info form.GoodInfo) error {
	var inv model.Inventory
	inv.Goods = info.GoodID
	inv.Stocks = info.Num
	result := dao.Group.DB.Create(inv)
	return result.Error
}

func (i *inventory) GetInvDetails(GoodID int32) (*model.Inventory, error) {
	var inv model.Inventory
	if result := dao.Group.DB.Where("goods=?", GoodID).Find(&inv); result.Error != nil {
		return nil, result.Error
	}
	return &inv, nil
}
