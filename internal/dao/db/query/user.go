package query

import (
	"gorm.io/gorm"
	"mall/internal/dao"
	"mall/internal/form"
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
func (u *user) GetUserByMobile(mobile string) (model.User, error) {
	var user model.User
	result := u.Where("mobile=?", mobile).First(&user)
	return user, result.Error
}

func (u *user) GetUserByID(userID int64) (model.User, error) {
	var user model.User
	result := u.Where("id=?", userID).First(&user)
	return user, result.Error
}
func (u *user) CreateUser(user *model.User) error {
	result := u.Model(&model.User{}).Save(&user)
	return result.Error
}
func (u *user) GetAllUsers(pn, ps int64) ([]*model.User, error) {
	var users []*model.User
	result := u.Scopes(Paginate(pn, ps)).Find(&users)
	return users, result.Error
}

func (u *user) UpdateUserInfo(user2 form.UpdateUserInfo) error {
	updateMap := make(map[string]interface{})
	if user2.UserName != "" {
		updateMap["username"] = user2.UserName
	}
	if user2.BirthDay != "" {
		updateMap["birthday"] = user2.BirthDay
	}
	if user2.Avatar != "" {
		updateMap["avatar"] = user2.Avatar
	}

	result := u.Model(&model.User{}).Updates(updateMap)
	return result.Error
}

func (u *user) UpdateEmail(userID int64, email string) error {
	result := u.Model(&model.User{}).Where("id=?", userID).Update("email", email)
	return result.Error
}

func (u *user) ModifyPassword(userID int64, password string) error {
	result := u.Model(&model.User{}).Where("id=?", userID).Update("password", password)
	return result.Error
}
func (u *user) DeleteUser(userID int64) error {
	result := u.Where("id=?", userID).Delete(&model.User{})
	return result.Error
}

func (u *user) SearchUser(userName string) (users []model.User, num int64, err error) {
	result := u.Model(&model.User{}).Where("user_name like ?", "%"+userName+"%").Find(&users)
	return users, result.RowsAffected, result.Error
}
func (u *user) SearchUserByPage(userName string) (users []model.User, num int64, err error) {
	result := u.Scopes(Paginate(0, 0)).Model(&model.User{}).Where("user_name like ?", "%"+userName+"%").Find(&users)
	return users, result.RowsAffected, result.Error
}
