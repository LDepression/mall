package main

import (
	"context"
	"log"
	"mall/internal/global"
	"mall/internal/model/common"
	"mall/internal/routering/router"
	"mall/internal/setting"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}
func InitSetting() {
	setting.InitAll()
}
func main() {
	InitSetting()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", common.ValidateMobile) //这里的mobile和from表单里的是一样的
		_ = v.RegisterValidation("email", common.ValidateEmail)   //这里的email和from表单里的是一样的
		//翻译错误
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "非法的手机号码", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		_ = v.RegisterTranslation("email", global.Trans, func(ut ut.Translator) error {
			return ut.Add("email", "非法的邮箱", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("email", fe.Field())
			return t
		})
	}
	r := router.NewRouter()
	s := &http.Server{
		Addr:           global.Setting.Serve.Addr,
		Handler:        r,
		ReadTimeout:    global.Setting.Serve.ReadTimeout,
		WriteTimeout:   global.Setting.Serve.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//zap.S().Info(global.Setting.Captcha)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()
	//c, _ := rocketmq.NewPushConsumer(
	//	consumer.WithNameServer([]string{"192.168.28.16:9876"}),
	//	consumer.WithGroupName("mall-reback"),
	//)
	//if err := c.Subscribe("order_back", consumer.MessageSelector{}, logic.AutoReback); err != nil {
	//	fmt.Println("读取消息失败")
	//}
	//_ = c.Start()
	//time.Sleep(time.Hour) //不能让主携程退出

	//c1, _ := rocketmq.NewPushConsumer(
	//	consumer.WithNameServer([]string{"192.168.28.16:9876"}),
	//	consumer.WithGroupName("mall-timeout"),
	//)
	//if err := c1.Subscribe("order_timeout", consumer.MessageSelector{}, logic.OrderTimeout); err != nil {
	//	fmt.Println("读取消息失败")
	//}
	//_ = c1.Start()
	//_ = c1.Shutdown()
	// 退出通知
	quit := make(chan os.Signal, 1)
	// 等待退出通知
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//_ = c.Shutdown()
	log.Println("shutDone....")
	// 给几秒完成剩余任务
	ctx, cancel := context.WithTimeout(context.Background(), global.Setting.Serve.DefaultTimeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil { // 优雅退出
		log.Println("Server forced to ShutDown,Err:" + err.Error())
	}
}

// 优雅退出
//func gracefulExit(s *http.Server) {
//	// 退出通知
//	quit := make(chan os.Signal, 1)
//	// 等待退出通知
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//	log.Println("shutDone....")
//	// 给几秒完成剩余任务
//	ctx, cancel := context.WithTimeout(context.Background(), global.Setting.Serve.DefaultTimeout)
//	defer cancel()
//	if err := s.Shutdown(ctx); err != nil { // 优雅退出
//		log.Println("Server forced to ShutDown,Err:" + err.Error())
//	}
//}
