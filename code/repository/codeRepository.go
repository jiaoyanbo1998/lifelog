package repository

import (
	"context"
	"lifelog-grpc/code/repository/cache"
)

var (
	ErrCodeSendFrequent = cache.ErrCodeSendFrequent
	ErrCodeSendMany     = cache.ErrCodeSendMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) error
	SetBlackPhone(ctx context.Context, phone string) error
	IsBackPhone(ctx context.Context, phone string) (bool, error)
}

type CodeRepositoryV1 struct {
	codeCache cache.CodeCache
}

func NewCodeRepository(codeCache cache.CodeCache) CodeRepository {
	return &CodeRepositoryV1{
		codeCache: codeCache,
	}
}

func (c *CodeRepositoryV1) IsBackPhone(ctx context.Context, phone string) (bool, error) {
	return c.codeCache.IsBlackPhone(ctx, phone)
}

func (c *CodeRepositoryV1) SetBlackPhone(ctx context.Context, phone string) error {
	return c.codeCache.SetBlackPhone(ctx, phone)
}

func (c *CodeRepositoryV1) Set(ctx context.Context, biz string, phone string, code string) error {
	// 将验证码存储到redis中
	err := c.codeCache.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRepositoryV1) Verify(ctx context.Context, biz string, phone string, code string) error {
	// 将验证码存储到redis中
	err := c.codeCache.Verify(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}
