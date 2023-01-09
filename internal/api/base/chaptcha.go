package base

import (
	"mall/internal/global"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var Store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(global.Setting.Captcha.Width,
		global.Setting.Captcha.Height, global.Setting.Captcha.Length, global.Setting.Captcha.MaxSkew,
		global.Setting.Captcha.DotCount)
	cp := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorw("生产验证码错误", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
