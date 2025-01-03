package decorator

import (
	"context"
	"errors"
	"fmt"
	"lifelog-grpc/pkg/limitx"
	"lifelog-grpc/third/sms"
)

var (
	ErrLimited = errors.New("触发限流")
)

type SmsService struct {
	sendSmsService sms.SendSmsService
	limit          limitx.Limiter
}

func NewSmsService(sendSmsService sms.SendSmsService, limit limitx.Limiter) *SmsService {
	return &SmsService{
		sendSmsService: sendSmsService,
		limit:          limit,
	}
}

func (s *SmsService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	// key为限流对象，只要是唯一标识即可
	// ok == true 表示被限流了，不允许发送短信
	// ok == false 表示没有被限流，允许发送短信
	ok, err := s.limit.Limit(ctx, biz)
	if err != nil {
		return fmt.Errorf("短信服务限流异常 %w", err)
	}
	if ok == true {
		return ErrLimited
	}
	return s.sendSmsService.Send(ctx, biz, args, numbers...)
}
