package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// 该方法用于针对上传后的文件名格式化，简单来讲，将文件名 MD5 后再进行写入，防止直接把原始名称就暴露出去了

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
