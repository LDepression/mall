package logic

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mall/internal/dao/db/query"
	"mall/internal/form"
	"mall/internal/global"
	"mall/internal/middleware"
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

func UsingPayloadToFindUser(ctx *gin.Context) (user model.User, err errcode.Err) {
	payLoad, err := middleware.GetPayload(ctx)
	if err != nil {
		return user, errcode.ErrNotLogin
	}
	userID := payLoad.PalLoad.UserID
	quser := query.NewUser()
	user, err1 := quser.GetUserByID(userID)
	if err1 != nil {
		zap.S().Info("getUserByID failed,err:", err1)
		return user, errcode.ErrServer.WithDetails(err1.Error())
	}
	return user, nil
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
func (*user) Login(login *form.Login) (*reply.User, errcode.Err) {
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

func (user) Register(register *form.Register) errcode.Err {
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

func (user) GetUsers(pn, ps int64) ([]*reply.User, errcode.Err) {
	quser := query.NewUser()
	users, err := quser.GetAllUsers(pn, ps)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Info("目前暂无用户")
			return nil, errcode.ErrNotFound.WithDetails("目前没有用户")
		}
		return nil, errcode.ErrServer
	}
	usersInfo := make([]*reply.User, len(users))
	for i := range users {
		usersInfo = append(usersInfo, &reply.User{
			UserID:   int64(users[i].ID),
			UserName: users[i].UserName,
			Avatar:   users[i].Avatar,
			Email:    users[i].Email,
		})
	}
	return usersInfo, nil
}

func (user) GetUserByID(id int64) (reply.UserInfo2Visitor, errcode.Err) {
	var user reply.UserInfo2Visitor
	quser := query.NewUser()
	userInfo, err := quser.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errcode.ErrNotFound
		}
		return user, errcode.ErrServer
	}
	user.UserID = int64(userInfo.ID)
	user.Email = userInfo.Email
	user.Avatar = userInfo.Avatar
	user.UserName = userInfo.UserName
	user.Birthday = userInfo.BirthDay
	return user, nil
}

func (user) RefreshToken(req form.ReqRefresh) (string, errcode.Err) {
	refreshPayload, err := global.Maker.VerifyToken(req.RefreshToken)
	if err != nil {
		return "", errcode.ErrServer.WithDetails(err.Error())
	}
	if ok := refreshPayload.ExpiredAt.Before(time.Now()); ok {
		//此时已经过期了,就要重新去登录
		return "", errcode.ErrLoginTimeOut
	} else {
		newToken, _, err := global.Maker.CreateToken(refreshPayload.UserID, refreshPayload.UserName, global.Setting.Token.AccessTokenExpire)
		if err != nil {
			return "", errcode.ErrServer.WithDetails(err.Error())
		}
		return newToken, nil
	}
}

func (user) UpdateUserInfo(req form.UpdateUserInfo) errcode.Err {
	quser := query.NewUser()
	if _, err := quser.GetUserByID(req.UserID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	if err := quser.UpdateUserInfo(req); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (user) UpdateEmail(ctx *gin.Context, req form.UpdateEmail) errcode.Err {
	user, err := UsingPayloadToFindUser(ctx)
	if err != nil {
		return err
	}
	if ok := Group.Email.CheckCode(req.Code, user.Email); !ok {
		return errcode.ErrServer
	}
	quser := query.NewUser()
	err1 := quser.UpdateEmail(int64(user.ID), req.Email)
	if err1 != nil {
		return errcode.ErrServer.WithDetails(err1.Error())
	}
	return nil
}

func (user) ModifyPassword(ctx *gin.Context, req form.ModifyPassword) errcode.Err {
	user, err := UsingPayloadToFindUser(ctx)
	if err != nil {
		return err
	}
	if ok := Group.Email.CheckCode(req.Code, user.Email); !ok {
		return errcode.ErrServer
	}
	quser := query.NewUser()
	//将密码进行加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.NewPassword, options)
	ResPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	err1 := quser.ModifyPassword(req.UserID, ResPassword)
	if err1 != nil {
		return errcode.ErrServer
	}
	return nil
}
func (user) DeleteUser(userID int64) errcode.Err {
	quser := query.NewUser()
	if err := quser.DeleteUser(userID); err != nil {
		return errcode.ErrServer.WithDetails(err.Error())
	}
	return nil
}

func (user) SearchUser(username string) (users []reply.UserInfo2Visitor, num int64, err errcode.Err) {
	quser := query.NewUser()
	_, Total, err1 := quser.SearchUser(username)
	if err != nil {
		return users, Total, errcode.ErrServer.WithDetails(err1.Error())
	}
	num = Total
	usersInfo, _, err1 := quser.SearchUserByPage(username)
	if err1 != nil {
		return users, num, errcode.ErrServer.WithDetails(err1.Error())
	}
	for _, userInfo := range usersInfo {
		users = append(users, reply.UserInfo2Visitor{
			UserID:   int64(userInfo.ID),
			UserName: userInfo.UserName,
			Avatar:   userInfo.Avatar,
			Birthday: userInfo.BirthDay,
			Email:    userInfo.Email,
		})
	}
	return users, num, nil

}
