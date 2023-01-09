package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
	"mall/internal/pkg/app/errcode"
)

type user struct {
}

func NewUser() *user {
	return &user{}
}

func (*user) Login(ctx *gin.Context) {
	reqLogin := form.Login{}
	rly := app.NewResponse(ctx)
	err := ctx.ShouldBindJSON(&reqLogin)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	ok := base.Store.Verify(reqLogin.CaptchaID, reqLogin.CaptchaBS64, false)
	if !ok {
		zap.S().Error("base.Store.Verify failed")
		rly.Reply(errcode.ErrCaptcha.WithDetails("base.Store.Verify failed"), nil)
		return
	}
	////将验证码存入redis中去
	//dao.Group.Redis.Set(context.Background(), reqLogin.Mobile, reqLogin.CaptchaBS64, time.Duration(global.Setting.Captcha.Expired)*time.Second)
	userInfo, err := logic.Group.User.Login(ctx, &reqLogin)
	if err != nil {
		zap.S().Errorf("logic.Group.User.Login %s", err)
		rly.Reply(errcode.ErrServer, nil)
		return
	}
	rly.Reply(nil, userInfo)
}

func (*user) Register(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	register := form.Register{}
	err := ctx.ShouldBind(&register)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.User.Register(ctx, &register); err != nil {
		zap.S().Errorf(" logic.Group.User.Register failed,err:%v", err)
		rly.Reply(err, nil)
		return
	}
	rly.Reply(nil, nil)
	return
}

//func GetUsers(ctx *gin.Context) {
//	rly := app.NewResponse(ctx)
//
//	pn := ctx.DefaultQuery("pn", 1)
//	var users []reply.User
//	if users, err := logic.Group.User.GetUsers(); err != nil {
//		zap.S().Errorf("logic.Group.User.GetUsers() failed,err:%v", err)
//		rly.Reply(err, nil)
//		return
//	}
//	rly.Reply(nil, users)
//	return
//}
