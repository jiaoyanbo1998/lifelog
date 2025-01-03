package decorator

import (
	"context"
	"sync/atomic"
	"lifelog-grpc/third/sms"
)

type SmsServiceTimeout struct {
	sendSmsService []sms.SendSmsService // 服务商列表
	count          int32                // 连续超时的次数
	threshold      int32                // 阈值，连续超时次数超过阈值，就要切换
	index          int32                // 轮询的索引
}

func NewSmsServiceTimeout(sendSmsService []sms.SendSmsService, threshold int32) sms.SendSmsService {
	return &SmsServiceTimeout{
		sendSmsService: sendSmsService,
		threshold:      threshold,
	}
}

func (s *SmsServiceTimeout) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	// 获取当前服务商索引
	index := atomic.LoadInt32(&s.index)
	// 获取当前连续超时次数
	count := atomic.LoadInt32(&s.count)
	// 超过阈值，就切换
	if count > s.threshold {
		// 计算下一个服务商的下标
		newIndex := (index + 1) % int32(len(s.sendSmsService))
		// 切换
		if ok := atomic.CompareAndSwapInt32(&s.index, index, newIndex); ok {
			// 切换成功，重置超时计数
			atomic.StoreInt32(&s.count, 0)
		} // else {} 发生并发，其他协程已经切换成功
		// 更新当前服务商索引
		index = newIndex
	}
	// 获取当前索引的服务商
	smsService := s.sendSmsService[index]
	err := smsService.Send(ctx, biz, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		// 发生超时
		atomic.AddInt32(&s.count, 1)
		return err
	case nil:
		// 短信发送成功
		atomic.StoreInt32(&s.count, 0)
		return nil
	default:
		// 其他错误
		return err
	}
}
