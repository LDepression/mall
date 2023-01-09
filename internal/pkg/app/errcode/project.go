package errcode

var (
	UserNotFound   = NewErr(20001, "用户未查询到")
	UserExist      = NewErr(20002, "用户已存在")
	ErrPassword    = NewErr(20003, "用户名或密码错误")
	ErrSendTooMany = NewErr(20004, "发送次数过多")
	ErrCode        = NewErr(20005, "验证码有误或超时")
)
