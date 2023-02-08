package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func main() {
	//todo 支付宝的支付连接
	client, err1 := alipay.New("2021000122617027", "***", false)
	if err1 != nil {
		return
	}
	err1 = client.LoadAliPayPublicKey("***")
	if err1 != nil {
		return
	}
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://www.lyc666.xyz/pay/alipay/notify"
	p.ReturnURL = "http://127.0.0.1:8090/return"
	p.Subject = "1231"
	p.OutTradeNo = "2141213"
	p.TotalAmount = "123"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err1 := client.TradePagePay(p)
	if err1 != nil {
		zap.S().Info("生成支付url失败")
		return
	}
	fmt.Println(url.String())
}
