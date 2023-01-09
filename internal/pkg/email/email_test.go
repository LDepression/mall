package email

import (
	"fmt"
	"testing"
	"time"

	"github.com/0RAJA/Rutils/pkg/times"
)

var (
	Host     = "smtp.qq.com"
	Port     = 465
	UserName = "1197285120@qq.com"
	Password = "zibwqaowvgongjca"
	IsSSL    = true
	From     = "1197285120@qq.com"
	To       = []string{"1197285120@qq.com"}
)

/*
Host: smtp.qq.com
  Port: 465
  UserName: XXX@qq.com
  Password:
  IsSSL: true
  From: XXX@qq.com
  To:
    - XXX@qq.com
*/
func TestEmail_SendMail(t *testing.T) {
	defailtMailer := NewEmail(&SMTPInfo{
		Host:     Host,
		Port:     Port,
		IsSSL:    IsSSL,
		UserName: UserName,
		Password: Password,
		From:     From,
	})
	err := defailtMailer.SendMail( // 短信通知
		To,
		fmt.Sprintf("异常抛出，发生时间: %s,%d", times.GetNowDateTimeStr(), time.Now().Unix()),
		fmt.Sprintf("错误信息: %v", "NO"),
	)
	if err != nil {
		fmt.Println(err)
	}
}
