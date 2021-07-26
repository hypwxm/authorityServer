package aliyunOss

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var OssConfigDev = map[string]string{
	"host":     "https://authorityServerning-test.oss-cn-hangzhou.aliyuncs.com",
	"bucket":   "authorityServerning-test",
	"endpoint": "oss-cn-hangzhou.aliyuncs.com",
}

var OssConfigProd = map[string]string{
	"host":     "https://authorityServer.oss-cn-hangzhou.aliyuncs.com",
	"bucket":   "authorityServer",
	"endpoint": "oss-cn-hangzhou-internal.aliyuncs.com",
}

var OssConfig = make(map[string]string)

func init() {
	OssConfig = OssConfigProd
}

// 创建oss客户端
func CreateClient() *oss.Client {
	client, err := oss.New(OssConfig["endpoint"], "", "")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return client
}
