package sms

import "context"

type SendSmsService interface {
	Send(ctx context.Context, biz string, args []string, numbers ...string) error
}
