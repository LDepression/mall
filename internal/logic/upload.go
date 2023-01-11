package logic

import (
	"mall/internal/pkg/app/errcode"
	upload2 "mall/internal/upload"
	"mime/multipart"
)

type upload struct {
}

func (upload) UploadFile(file *multipart.FileHeader) (string, errcode.Err) {
	OSSClient := upload2.NewOss()
	url, err := OSSClient.UploadFile(file)
	if err != nil {
		return "", errcode.ErrServer.WithDetails(err.Error())
	}
	return url, nil
}
