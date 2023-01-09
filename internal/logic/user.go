package logic

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/global"
	"mall/internal/model"
	"mall/internal/model/reply"
	"mall/internal/pkg/app/errcode"
	"mall/internal/pkg/token"
	"mall/internal/work/email"
	"strings"
	"time"
)

type user struct{}

type tokenResult struct {
	token   string
	payload *token.Payload
	err     error
}

func createToken(resultChan chan<- tokenResult, userID int64, userName string, expireTime time.Duration) func() {
	return func() {
		defer close(resultChan)
		accessToken, pal, err := global.Maker.CreateToken(userID, userName, expireTime)
		resultChan <- tokenResult{
			token:   accessToken,
			payload: pal,
			err:     err,
		}
	}
}
func (*user) Login(ctx *gin.Context, login *form.Login) (*reply.User, errcode.Err) {
	//先去查询用户是否存在
	quser := query.NewUser()
	user, err := quser.GetUserByMobile(login.Mobile)
	if err != nil {
		return nil, errcode.UserNotFound
	}
	//验证密码
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(user.Password, "$")
	check := password.Verify(login.Password, passwordInfo[2], passwordInfo[3], options)
	if !check {
		return nil, errcode.ErrPassword
	}
	accessChan := make(chan tokenResult, 1)
	refreshChan := make(chan tokenResult, 1)
	//这里就是生成token
	global.Worker.SendTask(createToken(accessChan, int64(user.ID), user.UserName, global.Setting.Token.AccessTokenExpire))
	global.Worker.SendTask(createToken(refreshChan, int64(user.ID), user.UserName, global.Setting.Token.RefreshTokenExpire))
	accessResult := <-accessChan
	if accessResult.err != nil {
		zap.S().Error(accessResult.err)
		return nil, errcode.ErrServer
	}
	refreshResult := <-refreshChan
	if refreshResult.err != nil {
		zap.S().Error(refreshResult.err)
		return nil, errcode.ErrServer
	}
	//现在直接返回user就好了
	return &reply.User{
		UserID:       int64(user.ID),
		UserName:     user.UserName,
		Avatar:       user.Avatar,
		AccessToken:  accessResult.token,
		RefreshToken: refreshResult.token,
		Payload:      *accessResult.payload,
	}, nil
}

func (user) Register(ctx *gin.Context, register *form.Register) errcode.Err {
	User := model.User{}
	User.Mobile = register.Mobile
	User.Password = register.Password
	User.UserName = register.UserName
	User.Email = register.Email
	if ok := email.Check(register.EmailCode, register.Email); !ok {
		return errcode.ErrCode
	}
	quser := query.NewUser()
	if _, err := quser.GetUserByMobile(register.Mobile); err == gorm.ErrRecordNotFound {
		//将密码进行加密
		options := &password.Options{16, 100, 32, sha512.New}
		salt, encodedPwd := password.Encode(register.Password, options)
		newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
		User.Password = newPassword
		if err = quser.CreateUser(&User); err != nil {
			return errcode.ErrServer
		}
		return nil
	} else {
		return errcode.UserExist
	}
}

func (user) GetUsers() {
	quser := query.NewUser()
	quser.GetAllUsers()
}
