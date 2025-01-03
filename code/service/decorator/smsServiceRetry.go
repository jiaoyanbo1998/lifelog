package decorator

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/third/sms"
)

type SmsServiceRetry struct {
	sendSmsService sms.SendSmsService
	// 基础时间间隔
	baseInterval time.Duration
	// 最大重试次数
	maxRetry int
	// 最大时间间隔
	maxInterval time.Duration
	// 指数退避算法的指数
	exponent float64
	logger   loggerx.Logger
}

func NewSmsServiceRetry(sendSmsService sms.SendSmsService, maxRetry int,
	maxInterval time.Duration, baseInterval time.Duration, l loggerx.Logger, // 指数退避算法的指数
	exponent float64) *SmsServiceRetry {
	return &SmsServiceRetry{
		sendSmsService: sendSmsService,
		baseInterval:   baseInterval,
		maxRetry:       maxRetry,
		maxInterval:    maxInterval,
		logger:         l,
		exponent:       exponent,
	}
}
func (s *SmsServiceRetry) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	// 添加重试策略
	var err error
	for i := 0; i < s.maxRetry; i++ {
		err = s.sendSmsService.Send(ctx, biz, args, numbers...)
		// 短信发送成功
		if err == nil {
			return nil
		}
		// 重试
		s.logger.Error("短信发送失败，正在重试", loggerx.Error(err))
		waitTime := s.calcWaitTime(i + 1)
		time.Sleep(waitTime)
	}
	s.logger.Error("短信发送失败，重试次数过多，放弃重试",
		loggerx.Int("maxRetry", s.maxRetry))
	return fmt.Errorf("重试失败，%w", err)
}

// 指数退避算法 + 随机生成数：计算等待重试的时间
func (s *SmsServiceRetry) calcWaitTime(retry int) time.Duration {
	// 创建随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 指数退避算法：每次重试的时间间隔，递增
	currentInterval := time.Duration(float64(s.baseInterval) * math.Pow(s.exponent, float64(retry)))
	if currentInterval > s.maxInterval {
		currentInterval = s.maxInterval
	}
	// 创建随机数
	//    baseInterval.Nanoseconds()：获取纳秒数(int64)
	// 	  r.Int63n(n)：生成[0,n)之间的随机整数，n为int64整数
	// 	  time.Duration(...)：将毫秒数(int64)，转为，time.Duration
	randomness := time.Duration(r.Int63n(currentInterval.Nanoseconds()))
	return currentInterval + randomness
}
