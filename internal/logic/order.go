package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"mall/internal/dao"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

type order struct {
}

type orderListener struct {
	Code        string
	Detail      string
	ID          int32
	OrderAmount float32
}

//本地事务
func (l *orderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	checkedGoodsInfo, err := query.Group.ShopCart.GetShopCartCheckedGoodsByUserID(orderInfo.User)
	if err != nil {
		zap.S().Info("GetShopCartCheckedGoodsByUserID failed,err:", err)
		l.Code = "invalid argument"
		l.Detail = "通过用户获取购物车选中商品失败"
		return primitive.RollbackMessageState
	}
	if len(checkedGoodsInfo) == 0 {
		l.Code = "invalid argument"
		l.Detail = "没有选中商品"
		return primitive.RollbackMessageState
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
	if err := Group.Inventory.Sell(form.SellInfo{GoodsInfo: goodsInvInfo, OrderSn: orderInfo.OrderSn}); err != nil {
		l.Code = "resource exhaust"
		l.Detail = "库存不足"
		zap.S().Info("库存扣减失败")
		return primitive.CommitMessageState
	}

	tx := dao.Group.DB.Begin()

	//生成订单
	orderInfo.OrderMount = totalPrice
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		l.Code = "internal error"
		l.Detail = "内部错误"
		return primitive.CommitMessageState
	}
	l.OrderAmount = totalPrice
	l.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		l.Code = "internal error"
		l.Detail = "内部错误"
		return primitive.CommitMessageState
	}

	//最后将购物车中所选择的商品删除
	if result := tx.Where(&model.ShoppingCart{Checked: true, User: orderInfo.User}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		l.Code = "internal error"
		l.Detail = "清除购物车失败"
		return primitive.CommitMessageState
	}
	//发送延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.28.16:9876"}), producer.WithInstanceName("lyc"))
	if err != nil {
		panic("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	msg = primitive.NewMessage("orderTimeout", msg.Body)
	msg.WithDelayTimeLevel(2)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Info("发送延时消息失败")
		l.Code = "internal error"
		l.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}
	//if err = p.Shutdown(); err != nil {
	//	panic("关闭producer失败")
	//}
	l.Code = "OK"
	tx.Commit()
	return primitive.RollbackMessageState
}
func (l *orderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	if result := dao.Group.DB.Where(&model.OrderInfo{OrderSn: GenerateOrderSn(orderInfo.User)}).Find(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
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
	var listener orderListener
	p, err := rocketmq.NewTransactionProducer(
		&listener,
		producer.WithNameServer([]string{"192.168.28.16:9876"}),
		producer.WithInstanceName("kinase"),
	)
	if err != nil {
		zap.S().Info("rocketmq.NewTransactionProducer", err)
		return nil, errcode.ErrServer.WithDetails("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	order := &model.OrderInfo{
		User:         req.UserID,
		OrderSn:      GenerateOrderSn(req.UserID),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}
	jsonString, _ := json.Marshal(&order)
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_back", jsonString))
	zap.S().Info(res.String())
	if err != nil {
		zap.S().Info("SendMessageInTransaction failed", err)
		return nil, errcode.ErrServer.WithDetails("发送消息失败")
	}
	if listener.Code != "OK" {
		return nil, errcode.ErrServer.WithDetails("新建订单失败")
	}
	return &reply.OrderInfoResponse{OrderSn: order.OrderSn, UserId: req.UserID, Total: listener.OrderAmount}, nil
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

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo model.OrderInfo
		if err := json.Unmarshal(msgs[i].Body, &orderInfo); err != nil {
			zap.S().Info("反序列化失败")
			return consumer.ConsumeSuccess, err
		}
		//去查询订单是否已经支付了

		db := Init()
		tx := db.Begin()
		var order model.OrderInfo

		if result := tx.Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).Find(&order); result.RowsAffected == 0 {
			zap.S().Info("没有查询到订单")
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {

			//模拟向order_back发送消息
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.28.16:9876"}), producer.WithGroupName("lyc"), producer.WithInstanceName("yyy"))
			if err != nil {
				panic("生成producer失败")
			}
			if err = p.Start(); err != nil {
				panic(err)
			}
			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_back", msgs[i].Body))
			if err != nil {
				fmt.Println("发送失败", err)
				return consumer.ConsumeRetryLater, err
			}
		}
		order.Status = "TRADE_CLOSED"
		if result := tx.Save(&order); result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Info("保存订单信息失败")
			return consumer.ConsumeSuccess, result.Error
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
