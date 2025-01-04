package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/pkg/ginx/middleware"
	"lifelog-grpc/pkg/limitx"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/web"
	"net/http"
	"time"
)

// InitGin 初始化gin服务
func InitGin(userHandler *web.UserHandler, middlewares []gin.HandlerFunc,
	lifeLogHandler *web.LifeLogHandler, collectHandler *web.CollectHandler,
	commentHandler *web.CommentHandler, codeHandler *web.CodeHandler,
	interactive *web.InteractiveHandler) *gin.Engine {
	// 创建默认的gin服务
	server := gin.Default()
	// 注册中间件
	// 中间件注册，必须在路由注册之前
	server.Use(middlewares...)
	// 注册路由
	userHandler.RegisterRoutes(server)
	lifeLogHandler.RegisterRoutes(server)
	collectHandler.RegisterRoutes(server)
	commentHandler.RegisterRoutes(server)
	codeHandler.RegisterRoutes(server)
	interactive.RegisterRoutes(server)
	return server
}

// InitMiddlewares 初始化中间件
func InitMiddlewares(logger loggerx.Logger, cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		InitJwtMiddleware(logger, cmd),
		InitCORS(),
		InitLimitMiddleware(cmd, logger),
	}
}

// InitJwtMiddleware jwt中间件
func InitJwtMiddleware(logger loggerx.Logger, cmd redis.Cmdable) gin.HandlerFunc {
	return middleware.NewJwtMiddlewareBuilder(logger, cmd).
		IgnorePath("/user/register_email_password").
		IgnorePath("/user/login_email_password").
		IgnorePath("/code/send_code").
		IgnorePath("/user/login_phone_code").
		Builder()
}

// InitLimitMiddleware 限流中间件
func InitLimitMiddleware(cmd redis.Cmdable, logger loggerx.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 对整个web服务器限流，1秒，最多1000个请求
		limiter := limitx.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000)
		ok, err := limiter.Limit(ctx, "webook")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, web.Result[string]{
				Code: 500,
				Msg:  "系统错误",
				Data: "error",
			})
			logger.Error("限流器创建失败", loggerx.Error(err))
			return
		}
		if ok == true {
			ctx.JSON(http.StatusBadRequest, web.Result[string]{
				Code: 400,
				Msg:  "服务器繁忙，请稍后再试",
				Data: "error",
			})
			logger.Error("触发限流，可能有人在攻击你的系统")
			return
		}
		ctx.Next()
	}
}

// InitCORS 跨域中间件
func InitCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 指定哪些跨域请求可以访问服务器资源
		AllowOrigins: []string{"http://localhost:3000"},
		// 指定允许哪些http请求方式
		AllowMethods: []string{"POST", "GET"},
		// 指定跨域请求可以携带哪些请求头信息
		AllowHeaders: []string{"Authorization", "content-type"},
		// 指定哪些响应头信息可以暴露给客户端
		ExposeHeaders: []string{"Content-Length", "jwt-short-token", "jwt-long-token"},
		// 指定跨域请求是否可以携带凭证
		AllowCredentials: true,
		// 指定预检请求的结果可以缓存多久
		MaxAge: time.Minute * 20,
	})
}
