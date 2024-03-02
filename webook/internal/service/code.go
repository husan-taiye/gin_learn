package service

import (
	"context"
	"fmt"
	"gin_learn/webook/internal/repository"
	"gin_learn/webook/internal/repository/cache"
	"gin_learn/webook/internal/service/sms"
	"go.uber.org/zap"
	"math/rand"
)

const codeTplId = "18775556"

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phoneNum string) error
	Verify(ctx context.Context, biz string, phoneNum string, inputCode string) (bool, error)
}

type RSTCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
	logger *zap.Logger
	tplId  string
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service, tplId string, l *zap.Logger) CodeService {
	return &RSTCodeService{
		repo:   repo,
		smsSvc: smsSvc,
		logger: l,
		tplId:  tplId,
	}
}

func (cs *RSTCodeService) Send(ctx context.Context, biz string, phoneNum string) error {
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

func (cs *RSTCodeService) Verify(ctx context.Context, biz string, phoneNum string, inputCode string) (bool, error) {
	return cs.repo.Verify(ctx, biz, phoneNum, inputCode)
}

func (cs *RSTCodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
