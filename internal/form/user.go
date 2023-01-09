package form

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
