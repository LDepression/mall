package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/model"
)

type user struct {
	*gorm.DB
}

func NewUser() *user {
	return &user{
		dao.Group.DB,
	}
}
func (user) GetUserByMobile(mobile string) (model.User, error) {
	var user model.User
	result := dao.Group.DB.Where("mobile=?", mobile).First(&user)
	return user, result.Error
}

func (user) GetUserByID(userID int32) (model.User, error) {
	var user model.User
	result := dao.Group.DB.Where("id=?", userID).First(&user)
	return user, result.Error
}
func (user) CreateUser(user *model.User) error {
	result := dao.Group.DB.Model(&model.User{}).Save(&user)
	return result.Error
}
func (user) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	result := dao.Group.DB.Model(&model.User{}).Find(&user{})
	return users, result.Error
}
