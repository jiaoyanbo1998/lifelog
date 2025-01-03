package ioc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"lifelog-grpc/pkg/redisx"
)

func InitRedis() redis.Cmdable {
	// 连接redis
	type Config struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	config := Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	// 将配置文件中，redis下的所有配置项绑定到结构体字段上
	err := viper.UnmarshalKey("redis", &config)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	// 初始化PrometheusRedisHook，监控Redis的缓存命中率
	opts := prometheus.SummaryOpts{
		Namespace: "jib",
		Subsystem: "webook",
		Name:      "redis_response_time_accuracy",
		Help:      "监控redis响应时间和缓存命中率",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 0.5 == 50%的请求，0.05 == 误差
			0.7:  0.02,  // 0.7 == 70%的请求，0.02 == 误差
			0.9:  0.01,  // 0.9 == 90%的请求，0.01 == 误差
			0.95: 0.005, // 0.95== 95%的请求，0.005 == 误差
		},
	}
	redisx.NewPrometheusRedisHook(opts)
	return client
}
