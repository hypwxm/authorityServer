package aliyunOss

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 上传文件流
func UploadFileStream(fd io.Reader, filename string) (string, error) {
	client := CreateClient()
	// 获取存储空间。
	bucket, err := client.Bucket(OssConfig["bucket"])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 签名直传。
	// signedURL, err := bucket.SignURL(filename, oss.HTTPPut, 60)
	// if err != nil {
	// 	return "", err
	// }
	filename = strings.ReplaceAll(filename, "\\", "/")
	err = bucket.PutObject(filename, fd)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(OssConfig["host"]+"/"+filepath.Join(filename), "\\", "/"), nil
}
