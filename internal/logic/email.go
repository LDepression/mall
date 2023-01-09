package logic

import (
	"go.uber.org/zap"
	"mall/internal/global"
	"mall/internal/pkg/app/errcode"
	"mall/internal/work/email"
)

type pemail struct {
}

func (pemail) SendEmailCode(emailStr string) errcode.Err {
	sendEmail := email.NewSendEmailCode(emailStr)
	if err := email.CheckEmailExist(emailStr); err != nil {
		return errcode.ErrSendTooMany
	}
	global.Worker.SendTask(sendEmail.Task())
	go func() {
		result := sendEmail.Value()
		if result.Err != nil {
			switch result.Err {
			case email.ErrRequestTooMany:
				zap.S().Info(errcode.ErrTooManyRequests.Error(), zap.String("email:", emailStr))
			default:
				zap.S().Error(result.Err.Error())
			}
		}
	}()
	return nil
}

func (pemail) CheckCode(code, emailStr string) bool {
	return email.Check(code, emailStr)
}
