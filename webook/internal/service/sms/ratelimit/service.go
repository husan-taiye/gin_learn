package ratelimit

import (
	"context"
	"fmt"
	"gin_learn/webook/internal/service/sms"
	"gin_learn/webook/pkg/ratelimit"
)

type Service struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &Service{
		svc:     svc,
		limiter: limiter,
	}
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 加代码
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务限流出问题: %w", err)
	}
	if limited {
		return fmt.Errorf("触发限流")
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	// 加代码
	return err
}
