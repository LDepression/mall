package main

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func main() {
	//todo 支付宝的支付连接
	client, err1 := alipay.New("2021000122617027", "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPG+3BjpBezn2ZV6UrS3XSdLGu5hhST7yugdW9FjccKfgdg4Wht4fCOVverUoL0pr32Azic/bWXqYid6GysKBPkU/RWCp6HvhnLvjDWA7Xrv5NYzSOhRRfCsplFKuDOThQCzBRs4AyI0k/Ew9YUgdGZbKKnfulXkoL7xfnPWhpAC5xyi7HVoE1iDcZqh+/18TNDz9mj5gSCk0SylQ0gruplbivTnR/Y9GBSk2HH+zoGH5ZKSG4WSBXNMXzMNTAqvv4egAetBQuIs0PJ23+TcnSCm9i7E0HqVfqY0gX7uhJlDdS0dL2H+No0KcC+wZdrScNh1JHZv5Ue66B/dU++kmfAgMBAAECggEALKsjAw9ksLLzMTHSNDlRhWc3LshTPx30XuPXuzV77iaLc2riAyAuF8mfi1m7iwUpqcKeAvD/UtooMQn2Rrgq0ashps5GM+gG0O4wZ4jM2TWd2rnkZbfUK/ZKRlK1Vjg+A336jwXgNcYdzro8R+0SqD6GBA5bxbowkpkGdP8N3/fbON3omVyywCkqM0pkpURGmL6paOLv4uaOW+uPenl+LOUJrY0Xc2EteQ2sCSs5wKQzjuetXCAufE8t3UzaLGN1Fjuf+Arj3IEF5bDGOrIKH9MFhmZu8mI+u5sTnHdP2Z9HLjZ651uYPc5opm8qaNRm9/bxTNjXvxId0u3mjJXDwQKBgQD2TAopZQOWgG4oDS8sJ4PvYM6YNPF4TOHvrFkxWxp8lejytHPJo/B0L+o4oxA8TLLcROK/AlLJ8py70UBg4czkIc8VQqvRTsNJshxjOUucjYJEUFaGroJEdyaM+xdSI4PXYzW/bXGAZi2rwDanbSRJ0ibfsiCkLdGD1i42FylsfQKBgQCUvzc+D0Kldy4nHwtr8Lgz+I0do2WMUwsdZOqJMZkyFO+rapm6kdtGkDWrFNPiHeAfyaWk10tiwFSiJ5ypN7OGWM+u1+RFLs4c3AIuxLRs60NSGb50pG0LRHW1HIFU3EnVASjTyUaR8qKADk7Skz16div4ZH2o2CZAO7uuCrpVSwKBgQC+VakMsEmDhyCZxwaLYsYsuW0uZsusog2AQHp1D+h6Gwd4eEd6rjxmLZkdx7YGQ2d9naZ04gDEm77PzjeoJxdFvXLhBTLuf6WfUAtsCp0KELl4vXUAg9+btVCPZoNxMIz0aHDizFsIVO46akJSRJ/khIkSGf/roJTnAx+XIXMbqQKBgCl1KhZ39mFb2Fc7Bdnt28lQazRpiDIWKzDkIaJfWo1k3G/wZCxl3rHKG8s1IOES5pa1gx9wiihZ5rzTQYzAY578PdZDgfHuW7AbedhDJu58m+TUHAsykNnlGNGDwmA+ja16h0CQBsVC1RvP4RQ7yZTKPvMaPxPCOtcITwTxJIIFAoGBAMbHScG04Qe0BXx4ymq9mp13ActEJjPOgGMAIDpYDuoQKokeGwtOEvyjaoM0Uxd6uwR7tnp6DQ1IYB1T4BLpzg9thKW0XS5JUuNwnK4Vy5mdzX9SGtuIpI0isZ/aOtJ2FBs0VazGexd+FD5zr/QypjoOgED9Qu4htUtHitxXNfIF", false)
	if err1 != nil {
		return
	}
	err1 = client.LoadAliPayPublicKey("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmH+LezESf8moVBrU6MYuQGfYprfvSNGhz+CquieUORTZwMIISTTQ+npBqu4Zbh259d6o76viqTyOFv42h7FlRkq8v8mANsRCkcgmXCxoKHly3w+ZKF6QJx7s1t133Sn2m457yUUAAItPqnbx1Pyn9s5d8JCSWrbHBzFZU1QD4iFtZU/xWBFmVzyqAs9+lt6S0PcUp9K/cGCIZOxtQ1HMjvBtEPVMzxAuzMis+Qqi04ayO/rD/Ga2Nr68NpfLKvaCZdK0UzuvVGPD8AS8R3E6ICPiC+3793kou2EAJ/gqm0V/P/dd/UKRU+BU8Yph54GBVsuQgJ99ogYvD6aeHKgbwQIDAQAB")
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
