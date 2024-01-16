package service

import (
	"context"
	"fmt"
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "18775556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
	tplId  string
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service, tplId string) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
		tplId:  tplId,
	}
}

func (cs *CodeService) Send(ctx context.Context, biz string, phoneNum string) error {
	// 生成验证码
	code := cs.generateCode()
	// 放进redis
	err := cs.repo.Store(ctx, biz, phoneNum, code)
	if err != nil {
		return err
	}
	// 发送出去短信
	err = cs.smsSvc.Send(ctx, cs.tplId, []string{code}, phoneNum)
	return err
}

func (cs *CodeService) Verify(ctx context.Context, biz string, phoneNum string, inputCode string) (bool, error) {
	return cs.repo.Verify(ctx, biz, phoneNum, inputCode)
}

func (cs *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
