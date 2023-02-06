package query

import (
	"mall/internal/dao"
	"mall/internal/form"
	"mall/internal/model"
)

type order struct {
}

func (o *order) GetOrdersByUserIDByPage(req form.OrderFilterRequest) ([]*model.OrderInfo, error) {
	var rep []*model.OrderInfo
	result := dao.Group.DB.Scopes(Paginate(req.Page, req.PagePerNum)).Where(&model.OrderInfo{User: int32(req.UserID)}).Find(&rep)
	return rep, result.Error
}

func (o *order) CalculateRecordsByUserID(userID int32) (int64, error) {
	var total int64
	result := dao.Group.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: userID}).Count(&total)
	return total, result.Error
}

func (o *order) GetOrderByUserIDAndOrderID(userID int64, orderID int32) (*model.OrderInfo, error) {
	ordersInfo := &model.OrderInfo{}
	result := dao.Group.DB.Where(&model.OrderInfo{User: int32(userID), BaseModel: model.BaseModel{ID: orderID}}).First(&ordersInfo)
	return ordersInfo, result.Error
}

func (o *order) GetOrderGoodsInfo(orderID int32) ([]*model.OrderGoods, error) {
	var orderGoods []*model.OrderGoods
	result := dao.Group.DB.Where(&model.OrderGoods{Order: orderID}).Find(&orderGoods)
	return orderGoods, result.Error
}

func (o order) UpdateStatus(req form.OrderStatus) error {
	var orderInfo model.OrderInfo
	orderInfo.Status = req.Status
	result := dao.Group.DB.Where(&model.OrderInfo{OrderSn: req.OrderSn}).Updates(&orderInfo)
	return result.Error
}
