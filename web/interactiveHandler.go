package web

import (
	"github.com/gin-gonic/gin"
	interactivev1 "lifelog-grpc/api/proto/gen/api/proto/interactive/v1"
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/job"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
)

type InteractiveHandler struct {
	logger                   loggerx.Logger
	interactiveServiceClient interactivev1.InteractiveServiceClient
	job                      job.Job
	biz                      string
	syncProducer             lifeLogEvent.Producer
	JWTHandler
}

func NewInteractiveHandler(
	l loggerx.Logger,
	interactiveServiceClient interactivev1.InteractiveServiceClient,
	job job.Job,
	syncProducer lifeLogEvent.Producer) *InteractiveHandler {
	return &InteractiveHandler{
		logger:                   l,
		interactiveServiceClient: interactiveServiceClient,
		biz:                      "lifeLog",
		job:                      job,
		syncProducer:             syncProducer,
	}
}

func (a *InteractiveHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/interactive")
	// 点赞LifeLog
	rg.PUT("/like/:id", a.Like)
	// 取消点赞LifeLog
	rg.PUT("/unlike/:id", a.UnLike)
	// 收藏LifeLog
	rg.POST("/collect", a.Collect)
	// 取消收藏LifeLog
	rg.PUT("/unCollect/:id", a.UnCollect)
}

// Like 点赞LifeLog
func (i *InteractiveHandler) Like(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Like"))
		return
	}
	userInfo, ok := i.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:Like"))
		return
	}
	_, err = i.interactiveServiceClient.Like(ctx.Request.Context(), &interactivev1.LikeRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:    i.biz,
			BizId:  lifeLogId,
			UserId: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "点赞失败",
			Data: "error",
		})
		i.logger.Error("点赞失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:Like"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "点赞成功",
		Data: "success",
	})
}

// UnLike 取消点赞LifeLog
func (i *InteractiveHandler) UnLike(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnLike"))
		return
	}
	userInfo, ok := i.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:UnLike"))
		return
	}
	// 减少点赞数
	_, err = i.interactiveServiceClient.UnLike(ctx.Request.Context(), &interactivev1.UnLikeRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:    i.biz,
			BizId:  lifeLogId,
			UserId: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "取消点赞失败",
			Data: "error",
		})
		i.logger.Error("取消点赞失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnLike"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "取消点赞成功",
		Data: "success",
	})
}

// Collect 收藏LifeLog
func (i *InteractiveHandler) Collect(ctx *gin.Context) {
	type CollectReq struct {
		BizId int64 `json:"biz_id"`
	}
	var req CollectReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:collect"))
		return
	}
	userInfo, ok := i.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:collect"))
		return
	}
	_, err = i.interactiveServiceClient.Collect(ctx.Request.Context(), &interactivev1.CollectRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:    i.biz,
			BizId:  req.BizId,
			UserId: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "收藏失败",
			Data: "error",
		})
		i.logger.Error("收藏失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:collect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "收藏成功",
		Data: "success",
	})
}

// UnCollect 取消收藏LifeLog
func (i *InteractiveHandler) UnCollect(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnCollect"))
		return
	}
	userInfo, ok := i.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:UnCollect"))
		return
	}
	// 减少收藏数
	_, err = i.interactiveServiceClient.UnCollect(ctx.Request.Context(), &interactivev1.UnCollectRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:    i.biz,
			BizId:  lifeLogId,
			UserId: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "取消收藏失败",
			Data: "error",
		})
		i.logger.Error("取消收藏失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnCollect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "取消收藏成功",
		Data: "success",
	})
}
