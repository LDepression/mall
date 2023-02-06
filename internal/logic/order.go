package logic

import (
	"fmt"
	"go.uber.org/zap"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"math/rand"
	"time"
)

type order struct {
}

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}
func (o *order) OrderList(req form.OrderFilterRequest) (*reply.OrderListResponse, errcode.Err) {
	response := reply.OrderListResponse{}
	var total int64
	total, err := query.Group.Order.CalculateRecordsByUserID(int32(req.UserID))
	if err != nil {
		zap.S().Info("CalculateRecordsByUserID ,err:", err)
		return nil, errcode.ErrServer
	}
	OrderInfos, err := query.Group.Order.GetOrdersByUserIDByPage(req)
	for _, orderInfo := range OrderInfos {
		response.Data = append(response.Data, &reply.OrderInfoResponse{
			Id:          orderInfo.ID,
			UserId:      orderInfo.User,
			Name:        orderInfo.SignerName,
			Address:     orderInfo.Address,
			Post:        orderInfo.Post,
			OrderSn:     orderInfo.OrderSn,
			OrderStatus: orderInfo.Status,
			AddTime:     orderInfo.BaseModel.CreatedAt.Format("2006-01-02 15:04:05"),
			OrderType:   orderInfo.PayType,
			Total:       orderInfo.OrderMount,
		})
	}
	response.Total = total
	return &response, nil
}
func (o *order) CreateOrder(req form.OrderRequest) (*reply.OrderInfoResponse, errcode.Err) {
	/*
		1.先去购物车中将checked的商品给查询出来------>ids   nums
		2.通过ids去查询各个价格,然后计算总价
		3.然后去扣减库存-----调用库存的sell方法
		4.生成订单号
		5.删除购物车已经购买了的信息
	*/
	checkedGoodsInfo, err := query.Group.ShopCart.GetShopCartCheckedGoodsByUserID(req.UserID)
	if err != nil {
		zap.S().Info("GetShopCartCheckedGoodsByUserID failed,err:", err)
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	if len(checkedGoodsInfo) == 0 {
		return nil, errcode.ErrNotFound
	}
	ids := make([]int32, 0)
	NumsMap := make(map[int32]int32)
	for _, goodInfo := range checkedGoodsInfo {
		ids = append(ids, goodInfo.Goods)
		NumsMap[goodInfo.Goods] = goodInfo.Nums
	}

	//调用商品服务
	var totalPrice float32
	goodsInfo, err := Group.Good.BatchGetGoods(ids)
	//注意这里要传指针进来,要不然其他变量对他赋值的时候修改不了
	orderGoods := make([]*model.OrderGoods, 0)
	goodsInvInfo := make([]*form.GoodInfo, 0)
	for _, goodInfo := range goodsInfo {
		totalPrice += float32(NumsMap[goodInfo.ID]) * goodInfo.ShopPrice
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goodInfo.ID,
			GoodsName:  goodInfo.Name,
			GoodsImage: goodInfo.GoodsFrontImage,
			GoodsPrice: goodInfo.ShopPrice,
			Nums:       NumsMap[goodInfo.ID],
		})
		goodsInvInfo = append(goodsInvInfo, &form.GoodInfo{
			GoodID: goodInfo.ID,
			Num:    NumsMap[goodInfo.ID],
		})
	}

	//调用库存服务
	if err := Group.Inventory.Sell(form.SellInfo{GoodsInfo: goodsInvInfo}); err != nil {
		zap.S().Info("库存扣减失败")
		return nil, err
	}
	tx := dao.Group.DB.Begin()

	//生成订单
	order := &model.OrderInfo{
		User:         req.UserID,
		OrderSn:      GenerateOrderSn(req.UserID),
		OrderMount:   totalPrice,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	if result := tx.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errcode.ErrServer.WithDetails("生成订单失败")
	}

	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}

	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errcode.ErrServer.WithDetails("生成订单失败")
	}

	//最后将购物车中所选择的商品删除
	if result := tx.Where(&model.ShoppingCart{Checked: true, User: req.UserID}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errcode.ErrServer.WithDetails("生成订单失败")
	}
	tx.Commit()
	return &reply.OrderInfoResponse{OrderSn: order.OrderSn, UserId: req.UserID}, nil
}
func (o *order) OrderDetails(req form.OrderRequest) (*reply.OrderDetailsResponse, errcode.Err) {
	rep := &reply.OrderDetailsResponse{}
	//这个是返回的是订单的基本信息
	orderInfo, err := query.Group.Order.GetOrderByUserIDAndOrderID(int64(req.UserID), req.ID)
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	//这里返回的是商品的基本信息
	GoodsInfos, err := query.Group.Order.GetOrderGoodsInfo(req.ID)
	if err != nil {
		return nil, errcode.ErrServer.WithDetails(err.Error())
	}
	rep.Total = int32(len(GoodsInfos))
	for _, goodInfo := range GoodsInfos {
		rep.GoodData = append(rep.GoodData, &reply.OrderItemResponse{
			GoodID:    goodInfo.Goods,
			GoodName:  goodInfo.GoodsName,
			GoodImage: goodInfo.GoodsImage,
			Price:     goodInfo.GoodsPrice,
			Num:       goodInfo.Nums,
		})
	}
	rep.OrderData = &reply.OrderInfoResponse{
		Id:          orderInfo.ID,
		UserId:      orderInfo.User,
		Name:        orderInfo.SignerName,
		Address:     orderInfo.Address,
		Post:        orderInfo.Post,
		OrderSn:     orderInfo.OrderSn,
		OrderStatus: orderInfo.Status,
		AddTime:     orderInfo.Address,
		OrderType:   orderInfo.PayType,
		Total:       orderInfo.OrderMount,
	}
	return rep, nil
}
func (o *order) UpdateOrderStatus(req form.OrderStatus) errcode.Err {
	if err := query.Group.Order.UpdateStatus(req); err != nil {
		zap.S().Info("UpdateStatus failed,err:", err)
		return errcode.ErrServer
	}
	return nil
}
