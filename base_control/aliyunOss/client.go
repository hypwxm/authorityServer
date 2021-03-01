package aliyunOss

import (
	"babygrow/config"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var OssConfigDev = map[string]string{
	"host":     "https://babygrow-test.oss-cn-hangzhou.aliyuncs.com",
	"bucket":   "babygrow-test",
	"endpoint": "oss-cn-hangzhou.aliyuncs.com",
}

var OssConfigProd = map[string]string{
	"host":     "https://babygrow.oss-cn-hangzhou.aliyuncs.com",
	"bucket":   "babygrow",
	"endpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
}

var OssConfig = make(map[string]string)

func init() {
	if config.Env == config.ENV_DEV {
		OssConfig = OssConfigDev
	}
	if config.Env == config.ENV_PROD {
		OssConfig = OssConfigProd
	}
}

// 创建oss客户端
func CreateClient() *oss.Client {
	client, err := oss.New(OssConfig["endpoint"], "LTAI4G1c4NhVtABdPUm3HNoK", "DYF3hzWYvLzLli8SupXqYHu0fOdL04")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return client
}
