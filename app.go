package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	// web服务器
	server *gin.Engine
	// 定时任务
	cron *cron.Cron
}
