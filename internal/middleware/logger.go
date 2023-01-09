package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const Body = "body"

// ErrLogMsg 日志数据
func ErrLogMsg(c *gin.Context) []zap.Field {
	var body string
	data, ok := c.Get(Body)
	if ok {
		body = string(data.([]byte))
	}
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	fields := []zap.Field{
		zap.Int("status", c.Writer.Status()),   // 状态码
		zap.String("method", c.Request.Method), // 请求方法
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("body", body)}
	return fields
}

// LogBody 读取body内容缓存下来，为之后打印日志做准备
func LogBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		_ = c.Request.Body.Close() //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Set("body", bodyBytes)
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.S().Info(path,
			zap.Int("status", c.Writer.Status()),   // 状态码
			zap.String("method", c.Request.Method), // 请求方法
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
