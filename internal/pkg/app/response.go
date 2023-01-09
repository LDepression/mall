package app

import (
	"net/http"

	"mall/internal/model/common"
	"mall/internal/pkg/app/errcode"

	"github.com/gin-gonic/gin"
)

type Response struct {
	c *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{c: ctx}
}

func (r *Response) Reply(err errcode.Err, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	if err == nil {
		err = errcode.StatusOK
	} else {
		data = nil
	}
	r.c.JSON(http.StatusOK, common.State{
		Code: err.ECode(),
		Msg:  err.Error(),
		Data: data,
	})
}

func (r *Response) ReplyList(err errcode.Err, datas ...interface{}) {
	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}
	if err == nil {
		err = errcode.StatusOK
	} else {
		data = nil
	}
	r.c.JSON(http.StatusOK, common.State{
		Code: err.ECode(),
		Msg:  err.Error(),
		Data: common.List{List: data},
	})
}
