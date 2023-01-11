package v1

import (
	"github.com/0RAJA/Rutils/pkg/utils"
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
	userInfo, err1 := logic.Group.User.Login(&reqLogin)
	if err1 != nil {
		zap.S().Errorf("logic.Group.User.Login %s", err)
		rly.Reply(err1, nil)
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
	if err := logic.Group.User.Register(&register); err != nil {
		zap.S().Errorf(" logic.Group.User.Register failed,err:%v", err)
		rly.Reply(err, nil)
		return
	}
	rly.Reply(nil, nil)
	return
}

func (*user) GetUsers(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	pNum := ctx.DefaultQuery("pn", "1")
	pSize := ctx.DefaultQuery("ps", "4")
	iNum := utils.StringToIDMust(pNum)
	iSize := utils.StringToIDMust(pSize)
	//有了gorm后面,这里默认或者不默认都没关系
	users, err := logic.Group.User.GetUsers(iNum, iSize)
	if err != nil {
		zap.S().Errorf("logic.Group.User.GetUsers() failed,err:%v", err)
		rly.Reply(err, nil)
		return
	}
	rly.ReplyList(nil, users)
}

func (*user) GetUserByID(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	IDStr := ctx.Param("id")
	i := utils.StringToIDMust(IDStr)
	userInfo, err := logic.Group.User.GetUserByID(i)
	if err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil, userInfo)
}
func (*user) RefreshToken(ctx *gin.Context) {
	//token没有过期的话,也是可以刷新的
	rly := app.NewResponse(ctx)
	reqRefresh := form.ReqRefresh{}
	err := ctx.ShouldBind(&reqRefresh)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	token, err1 := logic.Group.User.RefreshToken(reqRefresh)
	if err != nil {
		rly.Reply(err1)
		zap.S().Infof("logic.Group.User.RefreshToken(ctx) failed,err:%v", err1)
		return
	}
	rly.Reply(nil, token)
}

func (*user) UpdateUserInfo(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	UpdateUserInfo := form.UpdateUserInfo{}
	if err := ctx.ShouldBind(&UpdateUserInfo); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.User.UpdateUserInfo(UpdateUserInfo); err != nil {
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	rly.Reply(nil)
}
func (*user) UpdateEmail(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	UpdateEmail := form.UpdateEmail{}
	if err := ctx.ShouldBind(&UpdateEmail); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.User.UpdateEmail(ctx, UpdateEmail); err != nil {
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	rly.Reply(nil)
}

func (*user) ModifyPassword(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	ModifyPassword := form.ModifyPassword{}
	if err := ctx.ShouldBind(&ModifyPassword); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.User.ModifyPassword(ctx, ModifyPassword); err != nil {
		rly.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	rly.Reply(nil)
}

func (*user) Delete(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var deleteUser form.DeleteUser
	err := ctx.ShouldBindJSON(&deleteUser)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if err := logic.Group.User.DeleteUser(deleteUser.UserID); err != nil {
		rly.Reply(err)
		return
	}
	rly.Reply(nil)
}

func (*user) SearchUsers(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	var req form.SearchName
	err := ctx.ShouldBind(&req)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	users, num, err1 := logic.Group.User.SearchUser(req.Username)
	if err1 != nil {
		rly.Reply(err1)
		return
	}
	rly.ReplyList(nil, gin.H{
		"users": users,
		"total": num,
	})
}
