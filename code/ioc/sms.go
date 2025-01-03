package ioc

import (
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/code/service/decorator"
	"lifelog-grpc/pkg/limitx"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/third/sms"
	"time"
)

func InitSms(logger loggerx.Logger, cmd redis.Cmdable) sms.SendSmsService {
	// 随机轮询
	pollingSmsService := decorator.NewSmsServicePolling(
		[]sms.SendSmsService{
			sms.NewMemorySmsService(logger),
			// sms.NewTencentService(nil, "", ""),
		}, logger)
	// 创建限流器对象
	// 限流阈值：1秒，100次请求
	limiter := limitx.NewRedisSlidingWindowLimiter(cmd, time.Second, 100)
	limitSmsService := decorator.NewSmsService(pollingSmsService, limiter)
	// 创建基于指数退避算法+随机数的重试策略
	retrySmsService := decorator.NewSmsServiceRetry(
		limitSmsService, 3, time.Second*3, time.Second*1, logger, 1.5)
	return retrySmsService
}
