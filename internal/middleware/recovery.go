package middleware

import (
	"fmt"
	"mall/internal/global"
	"mall/internal/pkg/email"
	"mall/internal/pkg/times"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 异常捕获处理
// 自定义 Recovery

// Recovery recover掉项目可能出现的panic
func Recovery(stack bool) gin.HandlerFunc {
	defailtMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.Setting.SMTPInfo.Host,
		Port:     global.Setting.SMTPInfo.Port,
		IsSSL:    global.Setting.SMTPInfo.IsSSL,
		UserName: global.Setting.SMTPInfo.UserName,
		Password: global.Setting.SMTPInfo.Password,
		From:     global.Setting.SMTPInfo.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				var body string
				data, ok := c.Get(Body)
				if ok {
					body = string(data.([]byte))
				}
				err1 := defailtMailer.SendMail( // 短信通知
					global.Setting.SMTPInfo.To,
					fmt.Sprintf("异常抛出，发生时间: %v", time.Now().Format(times.LayoutDateTime)),
					fmt.Sprintf("错误信息: %s\n,请求信息:%s\n,请求body:%s\n,调用堆栈信息:%s", err, string(httpRequest), body, string(debug.Stack())),
				)
				if err1 != nil {
					zap.S().Error(fmt.Sprintf("mail.SendMail Error: %v", err1.Error()))
				}

				// Check for a broken connection
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					zap.S().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					zap.S().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("body", body),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.S().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("body", body),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
