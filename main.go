package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	// 初始化web服务器
	app := InitApp()
	// 初始化Prometheus
	initPrometheus()
	// 初始化viper
	initViperDevelopment()
	// 启动定时任务
	app.cron.Start()
	// 结束定时任务
	//    这一句其实并没有结束定时任务，只是发出了一个结束定时任务的信号，
	//	  此后就不会再调度下一个任务，但是已经调度的任务还会继续执行
	ctx := app.cron.Stop()
	// 超时控制，避免定时任务不及时结束
	timer := time.NewTimer(time.Minute * 10)
	defer timer.Stop()
	select {
	case <-ctx.Done(): // 定时任务结束，会向此只读管道中发送一个结束信号
	case <-timer.C: // 计时器触发时，会向此只读管道中发送当前的时间
	}
	// 启动web服务器
	app.server.Run(":8080")
}

// 初始化Prometheus
func initPrometheus() {
	// 开启一个协程，避免阻塞主程序
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

// 初始化viper
func initViperDevelopment() {
	// 配置文件的名字
	viper.SetConfigName("dev")
	// 配置文件的类型
	viper.SetConfigType("yaml")
	// 当前目录下的，config目录
	viper.AddConfigPath("config")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
