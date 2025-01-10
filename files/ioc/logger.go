package ioc

import (
	"lifelog-grpc/pkg/loggerx"
)

func InitLogger() loggerx.Logger {
	return loggerx.NewZapLogger()
}
