package oss

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var OssClient *dysmsapi.Client

func App() *dysmsapi.Client {
	if OssClient == nil {
		OssClient, _ = dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4FseXioSz7PMPTWabTKM", "vGzMEjku0zN2p10sMJacwIOuOW2CyQ")
	}
	return OssClient
}
