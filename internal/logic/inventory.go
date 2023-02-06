package logic

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/pkg/app/errcode"
)

type inventory struct {
}

//SetInv 设置库存
func (i *inventory) SetInv(req form.GoodInfo) errcode.Err {
	qInv := query.NewInventory()
	if err := qInv.CheckGoodID(req.GoodID); err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errcode.ErrNotFound
		}
		return errcode.ErrServer
	}
	var inv model.Inventory
	inv.Goods = req.GoodID
	inv.Stocks = req.Num
	if err := qInv.SetInv(req); err != nil {
		zap.S().Info("qInv.SetInv(req) failed,", err.Error())
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

//InvDetails 库存信息
func (i *inventory) InvDetails(req form.GoodInfo) (*form.GoodInfo, errcode.Err) {
	qInv := query.NewInventory()
	if err := qInv.CheckGoodID(req.GoodID); err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errcode.ErrNotFound
		}
		return nil, errcode.ErrServer
	}
	InvInfo, err := qInv.GetInvDetails(req.GoodID)
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	return &form.GoodInfo{
		GoodID: InvInfo.Goods,
		Num:    InvInfo.Stocks,
	}, nil
}

//Sell 欲扣减
func (i *inventory) Sell(req form.SellInfo) errcode.Err {
	tx := dao.Group.DB.Begin()
	for _, GoodInfo := range req.GoodsInfo {
		for {
			var inv model.Inventory
			//先去查询一下
			if result := tx.Where("goods=?", GoodInfo.GoodID).Find(&inv); result.RowsAffected == 0 {
				//说明此时没哟这个消息的记录,那么此时扣减库存就失败了
				tx.Rollback()
				return errcode.ErrServer.WithDetails(result.Error.Error())
			}
			//这里先去判断一下,商品的库存是否够用
			//这里表示的是如果此时库存不够
			if inv.Stocks < GoodInfo.Num {
				tx.Rollback()
				return errcode.ErrServer.WithDetails("商品的库存不足....")
			}
			inv.Stocks -= GoodInfo.Num //扣减库存
			//update inventory set version=version+1 and stocks =stocks -1 where version =version and goods=goods
			if result := tx.Model(&model.Inventory{}).Where("version = ? and goods =?", inv.Version, inv.Goods).Select("version", "stocks").Updates(&model.Inventory{
				Stocks:  inv.Stocks,
				Version: inv.Version + 1,
			}); result.RowsAffected == 0 {
				zap.S().Info("扣减库存失败...")
			} else {
				break
			}
		}
	}
	tx.Commit()
	return nil
}

//////Reback 库存归还
//func (i *inventory) Reback() errcode.Err {
//
//}
