package email

import (
	"errors"
	"fmt"
	"mall/internal/global"
	"mall/internal/pkg/email"
	"mall/internal/pkg/utils"
	"sync"
	"time"
)

var ErrRequestTooMany = errors.New("请求次数过多")

type MaskEmail struct {
	sync.RWMutex
	m map[string]struct{}
}

type Bind struct {
	sync.RWMutex
	m map[string]string
}

func NewMaskEmail() *MaskEmail {
	return &MaskEmail{
		RWMutex: sync.RWMutex{},
		m:       make(map[string]struct{}),
	}
}
func NewBind() *Bind {
	return &Bind{
		RWMutex: sync.RWMutex{},
		m:       make(map[string]string),
	}
}

var mask = NewMaskEmail() //绑定邮箱
var bind = NewBind()      //绑定用户

type Result struct {
	Code string
	Err  error
}

type SendEmailCode struct {
	Email string
	ma    chan Result
}

func NewSendEmailCode(email string) *SendEmailCode {
	return &SendEmailCode{
		Email: email,
		ma:    make(chan Result),
	}
}

func CheckEmailExist(email string) error {
	_, ok := mask.m[email]
	if ok {
		return ErrRequestTooMany
	}
	return nil
}
func (s *SendEmailCode) Task() func() {
	return func() {
		mask.RWMutex.Lock()
		_, ok := mask.m[s.Email]
		if ok {
			mask.RWMutex.RUnlock()
			s.ma <- Result{Err: ErrRequestTooMany}
			return
		}
		mask.m[s.Email] = struct{}{}
		mask.RWMutex.Unlock()
		SendInfo := email.NewEmail(&email.SMTPInfo{
			Host:     global.Setting.SMTPInfo.Host,
			Port:     global.Setting.SMTPInfo.Port,
			IsSSL:    global.Setting.SMTPInfo.IsSSL,
			UserName: global.Setting.SMTPInfo.UserName,
			Password: global.Setting.SMTPInfo.Password,
			From:     global.Setting.SMTPInfo.From,
		})
		code := utils.RandomString(5)
		err := SendInfo.SendMail([]string{s.Email}, fmt.Sprintf("%s:验证码:%s", "你的", code), `😘`)
		if err != nil {
			s.ma <- Result{Err: err}
			return
		}
		bind.Lock()
		bind.m[s.Email] = code
		bind.Unlock()
		s.AfterDelete()
		s.ma <- Result{
			Code: code,
			Err:  nil,
		}
		close(s.ma)
	}
}

//延时删除
func (s *SendEmailCode) AfterDelete() {
	time.AfterFunc(global.Setting.Auto.SendEmailTime, func() {
		mask.RWMutex.RLock()
		delete(mask.m, s.Email)
		mask.RWMutex.Unlock()
	})
	time.AfterFunc(global.Setting.Auto.CodeValidTime, func() {
		bind.RWMutex.RLock()
		delete(bind.m, s.Email)
		bind.RWMutex.Unlock()
	})
}

func (s *SendEmailCode) Value() Result {
	return <-s.ma
}

func Check(code, email string) bool {
	bind.RLock()
	defer bind.RUnlock()
	if bind.m[email] == code {
		delete(bind.m, email)
		return true
	} else {
		return false
	}
}
