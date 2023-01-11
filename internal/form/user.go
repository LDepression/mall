package form

import "mime/multipart"

type Login struct {
	Mobile      string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Password    string `json:"password" from:"password" binding:"required,min=3,max=10"`
	CaptchaID   string `json:"captchaID" form:"captchaID" binding:"required"`
	CaptchaBS64 string `json:"captchaBS64" form:"captchaBS64" binding:"required,min=5,max=5"`
}

type Register struct {
	UserName   string `json:"userName" form:"userName" binding:"required"`
	Mobile     string `json:"mobile" form:"mobile" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required"`
	RePassword string `json:"RePassword" form:"RePassword" binding:"required,eqfield=Password"`
	EmailCode  string `json:"emailCode" binding:"required"`
	Birthday   string `json:"birthday"`
	Email      string `json:"email" binding:"required"`
}

type ReqRefresh struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UpdateUserInfo struct {
	UserID   int64  `json:"userID"`
	Avatar   string `json:"avatar"`
	BirthDay string `json:"birthDay"`
	UserName string `json:"userName"`
}

type UpdateEmail struct {
	UserID int64  `json:"userID"`
	Code   string `json:"code"`
	Email  string `json:"email" binding:"required,email"`
}

type ModifyPassword struct {
	UserID      int64  `json:"userID" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
	Code        string `json:"code"`
}

type DeleteUser struct {
	UserID int64 `json:"userID" binding:"required"`
}

type SearchName struct {
	Username string `json:"username"`
}

type FileUpload struct {
	File *multipart.FileHeader `json:"file" binding:"required" form:"file"`
}
