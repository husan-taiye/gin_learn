package ali

import sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

type Service struct {
	templateCode *string
	signName     *string
	client       *sms.Client
}

func NewService(client *sms.Client, templateCode string, signName string) *Service {
	return &Service{
		client:       client,
		templateCode: &templateCode,
		signName:     &signName,
	}
}
