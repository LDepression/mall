package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

/*
Gomail 是一个用于发送电子邮件的简单又高效的第三方开源库，
目前只支持使用 SMTP 服务器发送电子邮件，但是其 API 较为灵活，如果有其它的定制需求也可以轻易地借助其实现，这恰恰好符合我们的需求，
因为目前我们只需要一个小而美的发送电子邮件的库就可以了。
*/

// SMTPInfo 传递发送邮箱所必需的信息
type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

type Email struct {
	*SMTPInfo
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()                                           // 创建一个消息实例，可以用于设置邮件的一些必要信息
	m.SetHeader("From", e.From)                                        // 发件人
	m.SetHeader("To", to...)                                           // 收件人
	m.SetHeader("Subject", subject)                                    // 主题
	m.SetBody("text/html", body)                                       // 正文
	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password) // 建一个新的 SMTP 拨号实例，设置对应的拨号信息用于连接 SMTP 服务器
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m) // 打开与 SMTP 服务器的连接并发送电子邮件
}
