package model

type User struct {
	BaseModel
	Mobile   string `gorm:"type:varchar(100);not null;index:idx_mobile;unique"`
	UserName string `gorm:"type:varchar(200);not null"`
	Password string `gorm:"type:varchar(200);not null"`
	Avatar   string `gorm:"type:varchar(200);default:'https://tupian.qqw21.com/article/UploadPic/2019-4/20194292022539844.jpeg'"`
	Gender   string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女 male表示男'"`
	Role     int    `gorm:"column:role;default:1; type:int comment '1表示普通用户,2表示管理员'"`
	Email    string `gorm:"type:varchar(200);not null"`
	BirthDay string `gorm:"column:birthday"`
}
