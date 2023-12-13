package ratelimit

import (
	"context"
	"fmt"
	"github.com/xuqil/webook/internal/service/sms"
	"github.com/xuqil/webook/pkg/ratelimit"
)

var ErrLimited = fmt.Errorf("触发了限流")

type SMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &SMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *SMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题, %w", err)
	}
	if limited {
		return ErrLimited
	}

	return s.svc.Send(ctx, tpl, args, numbers...)
}
