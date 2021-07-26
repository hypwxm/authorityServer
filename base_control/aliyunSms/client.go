package aliyunSms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var OssClient *dysmsapi.Client

func App() *dysmsapi.Client {
	if OssClient == nil {
		OssClient, _ = dysmsapi.NewClientWithAccessKey("cn-hangzhou", "", "")
	}
	return OssClient
}
