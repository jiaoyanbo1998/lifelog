package test

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	server := gin.Default()
	initPrometheus()
	server.POST("/hello", func(ctx *gin.Context) {
		var u User
		ctx.Bind(&u)
		r := rand.Int31n(1000)
		time.Sleep(time.Millisecond * time.Duration(r))
		// 这里我们模拟一下错误
		// 模拟 10% 比例的错误
		if r%100 < 10 {
			ctx.String(http.StatusInternalServerError, "系统错误")
		} else {
			ctx.String(http.StatusOK, u.Name)
		}
	})
	server.Run(":8080")
}

type User struct {
	Name string `json:"name"`
}

func initPrometheus() {
	// 开启一个协程，避免阻塞主程序
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
