package v1

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
)

func SendEmailCode(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	sendEmail := form.SendEmailCode{}
	err := ctx.ShouldBind(&sendEmail)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.Email.SendEmailCode(sendEmail.Email); err != nil {
		rly.Reply(err, nil)
		return
	}
	rly.Reply(nil)
}
