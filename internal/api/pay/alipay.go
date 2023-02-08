package pay

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"mall/internal/form"
	"mall/internal/global"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
)

//Notify 完成回调通知
func Notify(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	client, err1 := alipay.New(global.Setting.AliPay.AppID, global.Setting.AliPay.PrivateKey, global.Setting.AliPay.IsProduction)
	if err1 != nil {
		zap.S().Info("实例化支付宝失败")
		rly.Reply(errcode.ErrServer.WithDetails(err1.Error()))
		return
	}
	err1 = client.LoadAliPayPublicKey(global.Setting.AliPay.AliPublicKey)
	if err1 != nil {
		zap.S().Info("加载支付宝公钥失败")
		rly.Reply(errcode.ErrServer.WithDetails(err1.Error()))
		return
	}
	var noti, err = client.GetTradeNotification(ctx.Request)
	if err != nil {
		zap.S().Info("获取支付信息失败...")
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	if err := logic.Group.Order.UpdateOrderStatus(form.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	}); err != nil {
		zap.S().Info("logic.Group.Order.UpdateOrderStatus failed,err:", err)
		rly.Reply(errcode.ErrServer.WithDetails("更改订单状态失败"))
		return
	}
	alipay.AckNotification(ctx.Writer) // 确认收到通知消息
}
