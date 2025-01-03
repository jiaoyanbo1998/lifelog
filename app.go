package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"lifelog-grpc/event/lifeLogEvent"
)

type App struct {
	// web服务器
	server *gin.Engine
	// 消费者组
	consumers []lifeLogEvent.Consumer
	// 定时任务
	cron *cron.Cron
}
