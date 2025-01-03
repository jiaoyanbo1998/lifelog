package sms

import (
	"context"
	"lifelog-grpc/pkg/loggerx"
)

type MemorySmsService struct {
	logger loggerx.Logger
}

func NewMemorySmsService(l loggerx.Logger) *MemorySmsService {
	return &MemorySmsService{
		logger: l,
	}
}

func (m *MemorySmsService) Send(ctx context.Context, biz string, args []string,
	numbers ...string) error {
	m.logger.Info("获取短信验证码",
		loggerx.String("code：", args[0]),
		loggerx.String("phone：", numbers[0]))
	return nil
}
