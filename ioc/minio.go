package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"lifelog-grpc/pkg/miniox"
)

func InitMinio() *miniox.MinioClient {
	// 1.解析配置文件
	type config struct {
		Endpoint  string `yaml:"endpoint"`
		UseSSL    bool   `yaml:"useSSL"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
	}
	c := config{
		Endpoint:  "localhost:9000",
		UseSSL:    false,
		AccessKey: "RaEKllxArIYNXNu7WIay",
		SecretKey: "9SkofrJa1DA3vwF1NHEgbzt86ozoLr4b8rEtCJzS",
	}
	err := viper.UnmarshalKey("minio", &c)
	if err != nil {
		logx.Error("配置文件解析失败")
	}
	// 2.初始化minio
	client, err := miniox.NewMinioClient(c.Endpoint, c.AccessKey, c.SecretKey, c.UseSSL)
	if err != nil {
		logx.Error("初始化minio失败")
	}
	return client
}
