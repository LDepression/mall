package v1

import (
	"github.com/gin-gonic/gin"
	"mall/internal/api/base"
	"mall/internal/form"
	"mall/internal/logic"
	"mall/internal/pkg/app"
)

func UploadFile(ctx *gin.Context) {
	rly := app.NewResponse(ctx)
	f := form.FileUpload{}
	err := ctx.ShouldBind(&f)
	if err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	url, err1 := logic.Group.Upload.UploadFile(f.File)
	if err != nil {
		rly.Reply(err1, nil)
		return
	}
	rly.Reply(nil, url)
}
