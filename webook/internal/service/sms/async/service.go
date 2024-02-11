package async

import (
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/service/sms"
)

type SMSService struct {
	svc  sms.Service
	repo repository.SMSAsyncReqRepository
}

func NewSMSService() *SMSService {
	return &SMSService{}
}

//func (S SMSService) startAsync() {
//	go func() {
//		reqs, err := S.repo.Find()
//		for _, req := range reqs {
//
//		}
//	}()
//}
//
//func (S SMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
//	err := S.svc.Send(ctx, tpl, args, numbers...)
//	if err != nil {
//		S.repo.Store()
//	}
//}
