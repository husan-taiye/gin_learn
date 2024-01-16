package ali

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

type Service struct {
	templateCode *string
	signName     *string
	client       *dysmsapi.Client
}

func NewClient(accessKeyId string, accessKeySecret string) (*dysmsapi.Client, error) {
	config := &openapi.Config{}
	config.AccessKeyId = &accessKeyId
	config.AccessKeySecret = &accessKeySecret
	_result := &dysmsapi.Client{}
	_result, _err := dysmsapi.NewClient(config)
	return _result, _err
}

func NewService(client *dysmsapi.Client, templateCode string, signName string) *Service {
	return &Service{
		client:       client,
		templateCode: &templateCode,
		signName:     &signName,
	}
}

func (s Service) Send(ctx context.Context, phoneNumber []string, signName string,
	templateCode string, templateParams string) error {
	// 电话号码字符串拼接
	var phoneNumbers string
	phoneNumbers = strings.Join(phoneNumber, ",")
	sendReq := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  &phoneNumbers,
		SignName:      &signName,
		TemplateCode:  &templateCode,
		TemplateParam: &templateParams,
	}
	sendResp, err := s.client.SendSms(sendReq)
	if err != nil {
		return err
	}
	code := sendResp.Body.Code
	if !tea.BoolValue(util.EqualString(code, tea.String("OK"))) {
		fmt.Sprintln(tea.String("错误信息: " + tea.StringValue(sendResp.Body.Message)))
		return err
	}
	return nil
}
