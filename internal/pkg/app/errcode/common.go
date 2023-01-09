package errcode

// 以 1 开头表示公共错误码
var (
	StatusOK           = NewErr(0, "成功")
	ErrParamsNotValid  = NewErr(1001, "参数有误")
	ErrNotFound        = NewErr(1002, "未找到资源")
	ErrServer          = NewErr(1003, "系统错误")
	ErrTooManyRequests = NewErr(1004, "请求过多")
	ErrTimeOut         = NewErr(1005, "请求超时")
	ErrCaptcha         = NewErr(1006, "验证码错误或者是验证码过期了")
)
