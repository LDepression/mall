package upload

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mall/internal/global"
	"mime/multipart"
	"os"
	"path"
	"time"
)

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
func test() {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := "***"
	accessKeySecret := "***"
	bucketName := "lycmall2"
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	objectName := "mall/first.jpg"
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := `D:\360download\R.jpg`
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	}
}

type OssClient struct {
}

func (*OssClient) UploadFile(file *multipart.FileHeader) (string, error) {
	client, err := oss.New(global.Setting.OSS.Endpoint, global.Setting.OSS.AccessKeyId, global.Setting.OSS.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(global.Setting.OSS.BucketName)
	if err != nil {
		return "", err
	}
	objectName := global.Setting.OSS.BasePath + "/" + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	err = bucket.PutObject(objectName, f)
	if err != nil {
		return "", err
	}
	return global.Setting.OSS.BucketUrl + "/" + objectName, nil
}

func NewOss() OssServer {
	return &OssClient{}
}
