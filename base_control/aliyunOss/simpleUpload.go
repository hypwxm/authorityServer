package aliyunOss

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

	log.Println("filename:", filename)

	// 签名直传。
	// signedURL, err := bucket.SignURL(filename, oss.HTTPPut, 60)
	// if err != nil {
	// 	return "", err
	// }

	err = bucket.PutObject(filename, fd)
	if err != nil {
		return "", err
	}
	return filepath.Join(OssConfig["host"], filename), nil
}
