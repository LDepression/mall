package upload

import "mime/multipart"

type OssServer interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}
