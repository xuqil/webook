package retryable

import (
	"context"
	"errors"
	"github.com/xuqil/webook/internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryMax int
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	if err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试失败，超过最大次数")
}
