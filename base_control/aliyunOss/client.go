package aliyunOss

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const OssHost = "https://babygrowning-test.oss-cn-hangzhou.aliyuncs.com"

// 创建oss客户端
func CreateClient() *oss.Client {
	client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", "LTAI4G1c4NhVtABdPUm3HNoK", "DYF3hzWYvLzLli8SupXqYHu0fOdL04")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return client
}
