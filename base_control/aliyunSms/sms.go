package aliyunSms

import (
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func SendSms(phone string, code string) error {

	client := App()
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = "太子乐集团"
	request.TemplateCode = "SMS_176536119"
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	if response.Code != "OK" {
		fmt.Printf("%+v", response)
		return errors.New(response.Message)
	}
	return nil
}
