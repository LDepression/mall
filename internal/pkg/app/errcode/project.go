package errcode

var (
	UserNotFound    = NewErr(20001, "用户未查询到")
	UserExist       = NewErr(20002, "用户已存在")
	ErrPassword     = NewErr(20003, "用户名或密码错误")
	ErrSendTooMany  = NewErr(20004, "发送次数过多")
	ErrCode         = NewErr(20005, "验证码有误或超时")
	ErrLoginTimeOut = NewErr(20006, "登录超时,请重新登录")
	ErrGoodExsit    = NewErr(20007, "商品已经存在了")
)

var (
	ErrNotManager = NewErr(30001, "鉴权失败,权限不足")
	ErrNotLogin   = NewErr(30002, "鉴权失败,请先登录")
)

var (
	ErrRedis = NewErr(40001, "redis内部错误")
)
