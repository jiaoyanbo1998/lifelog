package service

import (
	"context"
	"fmt"
	"math/rand"
	"lifelog-grpc/code/repository"
	"lifelog-grpc/code/service/decorator"
	"lifelog-grpc/third/sms"
)

var (
	ErrCodeSendFrequent = repository.ErrCodeSendFrequent
	ErrCodeSendMany     = repository.ErrCodeSendMany
	ErrLimited          = decorator.ErrLimited
)

type CodeService interface {
	SendPhoneCode(ctx context.Context, phone, biz string) error
	VerifyPhoneCode(ctx context.Context, phone, code, biz string) error
	SetBlackPhone(ctx context.Context, phone string) error
	IsBackPhone(ctx context.Context, phone string) (bool, error)
}

type CodeServiceV1 struct {
	codeRepository   repository.CodeRepository
	memorySmsService sms.SendSmsService
}

func NewCodeService(codeRepository repository.CodeRepository, memorySmsService sms.SendSmsService) CodeService {
	return &CodeServiceV1{
		codeRepository:   codeRepository,
		memorySmsService: memorySmsService,
	}
}

func (c *CodeServiceV1) IsBackPhone(ctx context.Context, phone string) (bool, error) {
	return c.codeRepository.IsBackPhone(ctx, phone)
}

func (c *CodeServiceV1) SetBlackPhone(ctx context.Context, phone string) error {
	return c.codeRepository.SetBlackPhone(ctx, phone)
}

func (c *CodeServiceV1) VerifyPhoneCode(ctx context.Context, phone, code, biz string) error {
	// 将验证码存储到redis中
	err := c.codeRepository.Verify(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeServiceV1) SendPhoneCode(ctx context.Context, phone, biz string) error {
	// 获取6为验证码
	code := c.GetCode()
	// 将验证码存储到redis中
	err := c.codeRepository.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 调用短信服务
	err = c.memorySmsService.Send(ctx, biz, []string{code}, phone)
	if err != nil {
		return fmt.Errorf("发送短信验证码失败，%w", err)
	}
	return nil
}

// GetCode 随机生成6为的数字验证码
func (c *CodeServiceV1) GetCode() string {
	res := rand.Intn(1000000) // 生成[0,999999]/[0,1000000]之间的随机整数
	return fmt.Sprintf("%06d", res)
}
