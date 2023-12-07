package service

import (
	"context"
	"fmt"
	"github.com/xuqil/webook/internal/repository"
	"github.com/xuqil/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	//	生成一个验证码
	code := svc.generateCode()
	//	写入 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	//发送验证码
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
