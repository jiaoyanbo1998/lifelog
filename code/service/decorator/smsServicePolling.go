package decorator

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"time"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/third/sms"
)

var ErrAllSmsServiceFailed = errors.New("所有的短信服务商都失败")

// SmsServicePolling 轮询
type SmsServicePolling struct {
	smsServiceArr []sms.SendSmsService
	logger        loggerx.Logger
	pc            prometheus.Counter
}

func NewSmsServicePolling(
	smsServiceArr []sms.SendSmsService, l loggerx.Logger) sms.SendSmsService {
	// 创建计数器
	pc := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "jyb",
		Subsystem: "webook",
		Name:      "sms_all_failure_count",
		Help:      "所有的短信服务商都失败的次数",
	})
	// 将计数器注册到Prometheus
	prometheus.MustRegister(pc)
	return &SmsServicePolling{
		smsServiceArr: smsServiceArr,
		logger:        l,
		pc:            pc,
	}
}

func (s *SmsServicePolling) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	// 创建一个随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// r.Perm(len(切片))，将切片中的元素打乱顺序
	//	  names := []string{"张三", "李四", "王五"}
	//	  r.Perm(len(names)) ==> ["张三", "李四", "王五"]，["李四","张三", "王五"]，...
	randomIndex := r.Perm(len(s.smsServiceArr))
	for _, idx := range randomIndex {
		smsService := s.smsServiceArr[idx]
		err := smsService.Send(ctx, biz, args, numbers...)
		if err == nil {
			return nil
		}
		s.logger.Error("短信发送失败，切换下一个服务商", loggerx.Error(err))
	}
	s.logger.Error("所有服务商都发送失败，很有可能是自己的网络崩了")
	s.pc.Inc() // 增加计数器
	return ErrAllSmsServiceFailed
}
