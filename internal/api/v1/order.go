package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/middleware"
	"mall/internal/pkg/app"
	"strconv"
)

type order struct {
}

func NewOrder() *order {
	return &order{}
}

func (o *order) List(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var req form.OrderFilterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	payLoad, err := middleware.GetPayload(ctx)
	if err != nil {
		rly.Reply(err)
		return
	}
	userId := payLoad.PalLoad.UserID
	req.UserID = userId
	orderInfos, err := logic.Group.Order.OrderList(req)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.ReplyList(nil, orderInfos)

}
func (o *order) Create(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var orderRequest form.OrderRequest
	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	payLoad, err := middleware.GetPayload(ctx)
	if err != nil {
		rly.Reply(err)
		return
	}
	userId := payLoad.PalLoad.UserID
	orderRequest.UserID = int32(userId)
	rsp, err := logic.Group.Order.CreateOrder(orderRequest)
	if err != nil {
		zap.S().Info("logic.Group.Order.CreateOrder failed,err:", err)
		rly.Reply(err)
		return
	}
	reMap := make(map[string]interface{})
	reMap["order_sn"] = rsp.OrderSn
	reMap["user_id"] = rsp.UserId
	rly.Reply(nil, reMap)
}
func (o *order) Details(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	//这里的id是order的id
	idStr := ctx.Param("id")
	i, _ := strconv.Atoi(idStr)
	var req form.OrderRequest
	req.ID = int32(i)
	payLoad, err := middleware.GetPayload(ctx)
	if err != nil {
		rly.Reply(err)
		return
	}
	if payLoad.Role == 1 {
		req.UserID = int32(payLoad.PalLoad.UserID)
	}
	orderDetails, err := logic.Group.Order.OrderDetails(req)
	if err != nil {
		zap.S().Info("logic.Group.Order.OrderDetails(req) failed,err:", err)
		rly.Reply(err)
		return
	}

	reMap := make(map[string]map[string]interface{})
	for i, orderDetail := range orderDetails.GoodData {
		reMap[fmt.Sprintf("good_data%d", i+1)] = make(map[string]interface{})
		reMap[fmt.Sprintf("good_data%d", i+1)]["good_id"] = orderDetail.GoodID
		reMap[fmt.Sprintf("good_data%d", i+1)]["good_name"] = orderDetail.GoodName
		reMap[fmt.Sprintf("good_data%d", i+1)]["good_images"] = orderDetail.GoodImage
		reMap[fmt.Sprintf("good_data%d", i+1)]["good_num"] = orderDetail.Num
		reMap[fmt.Sprintf("good_data%d", i+1)]["good_price"] = orderDetail.Price
	}
	reMap["order_data"] = make(map[string]interface{})
	reMap["order_data"]["order_sn"] = orderDetails.OrderData.OrderSn
	reMap["order_data"]["order_address"] = orderDetails.OrderData.Address
	reMap["order_data"]["add_time"] = orderDetails.OrderData.AddTime
	reMap["order_data"]["name"] = orderDetails.OrderData.Name
	reMap["order_data"]["post"] = orderDetails.OrderData.Post
	reMap["order_data"]["user_id"] = orderDetails.OrderData.UserId
	rly.Reply(nil, reMap)
}
func (o *order) Update(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var updateStatus form.OrderStatus
	if err := ctx.ShouldBindJSON(&updateStatus); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	payLoad, err := middleware.GetPayload(ctx)
	if err != nil {
		rly.Reply(err)
		return
	}
	userId := payLoad.PalLoad.UserID
	updateStatus.UserID = userId
	if err := logic.Group.Order.UpdateOrderStatus(updateStatus); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil)
}
