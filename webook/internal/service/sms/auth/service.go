package auth

import (
	"context"
	"gin_learn/webook/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

// Send 发送，biz必须是线下申请的一个代表业务方的token
func (S SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var tc TokenClaims
	_, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return S.key, nil
	})
	if err != nil {
		return err
	}
	return S.svc.Send(ctx, biz, args, numbers...)
}

type TokenClaims struct {
	Tpl string
	jwt.RegisteredClaims
}
