package web

import (
	"github.com/gin-gonic/gin"
	interactivev1 "lifelog-grpc/api/proto/gen/interactive/v1"
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
	JWTHandler
}

func NewInteractiveHandler(
	l loggerx.Logger,
	interactiveServiceClient interactivev1.InteractiveServiceClient,
	job job.Job) *InteractiveHandler {
	return &InteractiveHandler{
		logger:                   l,
		interactiveServiceClient: interactiveServiceClient,
		biz:                      "lifeLog",
		job:                      job,
	}
}

func (a *InteractiveHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/interactive")
	// 点赞LifeLog
	rg.POST("/like", a.Like)
	// 取消点赞LifeLog
	rg.PUT("/unlike/:id", a.UnLike)
	// 收藏LifeLog
	rg.POST("/collect", a.Collect)
	// 取消收藏LifeLog
	rg.POST("/unCollect", a.UnCollect)
	// 关注作者
	rg.POST("/follow", a.Follow)
	// 取消关注作者
	rg.POST("/unfollow", a.UnFollow)
	// 关注列表
	rg.GET("/followList", a.FollowList)
	// 粉丝列表
	rg.GET("/fanList", a.FanList)
	// 互关列表
	rg.GET("/bothFollowList", a.BothFollowList)
}

// Like 点赞LifeLog
func (i *InteractiveHandler) Like(ctx *gin.Context) {
	type LikeReq struct {
		Id           int64 `json:"id"`
		TargetUserId int64 `json:"target_user_id"`
	}
	var req LikeReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		i.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:Like"))
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
			Biz:          i.biz,
			BizId:        req.Id,
			UserId:       userInfo.Id,
			TargetUserId: req.TargetUserId,
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
		BizId     int64 `json:"biz_id"`
		CollectId int64 `json:"collect_id"`
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
		}, CollectId: req.CollectId,
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
	type CollectReq struct {
		BizId     int64 `json:"biz_id"`
		CollectId int64 `json:"collect_id"`
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
			loggerx.String("method：", "InteractiveHandler:UnCollect"))
		return
	}
	// 减少收藏数
	_, err = i.interactiveServiceClient.UnCollect(ctx.Request.Context(), &interactivev1.UnCollectRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:    i.biz,
			BizId:  req.BizId,
			UserId: userInfo.Id,
		}, CollectId: req.CollectId,
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

func (a *InteractiveHandler) Follow(context *gin.Context) {
	// 获取用户信息
	info, ok := a.GetUserInfo(context)
	if !ok {
		context.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:UnFollow"))
		return
	}
	type FollowReq struct {
		FolloweeId int64 `json:"followee_id"`
	}
	var req FollowReq
	err := context.Bind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:Follow"))
		return
	}
	_, err = a.interactiveServiceClient.InsertFollow(context.Request.Context(),
		&interactivev1.InsertFollowRequest{
			Follow: &interactivev1.Follow{
				FollowerId: info.Id,
				FolloweeId: req.FolloweeId,
			},
		})
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "关注失败",
			Data: "error",
		})
		a.logger.Error("关注失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:Follow"))
		return
	}
	context.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "关注成功",
		Data: "success",
	})
}

func (a *InteractiveHandler) UnFollow(context *gin.Context) {
	// 获取用户信息
	info, ok := a.GetUserInfo(context)
	if !ok {
		context.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:UnFollow"))
		return
	}
	type FollowReq struct {
		FolloweeId int64 `json:"followee_id"`
	}
	var req FollowReq
	err := context.Bind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "请求参数错误",
			Data: "error",
		})
		a.logger.Error("请求参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnFollow"))
		return
	}
	_, err = a.interactiveServiceClient.CancelFollow(context.Request.Context(),
		&interactivev1.CancelFollowRequest{
			Follow: &interactivev1.Follow{
				FollowerId: info.Id,
				FolloweeId: req.FolloweeId,
			},
		})
	if err != nil {
		context.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "取消关注失败",
			Data: "error",
		})
		a.logger.Error("取消关注失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:UnFollow"))
		return
	}
	context.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "取消关注成功",
		Data: "success",
	})
}

// FollowList 关注列表
func (a *InteractiveHandler) FollowList(ctx *gin.Context) {
	// 获取用户信息
	info, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:FollowList"))
		return
	}
	list, err := a.interactiveServiceClient.FollowList(ctx.Request.Context(),
		&interactivev1.FollowListRequest{
			Id: info.Id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取关注列表失败",
			Data: "error",
		})
		a.logger.Error("获取关注列表失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:FollowList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]int64]{
		Code: 200,
		Msg:  "获取关注列表成功",
		Data: list.Ids,
	})
}

func (a *InteractiveHandler) BothFollowList(ctx *gin.Context) {
	// 获取用户信息
	info, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:BothFollowList"))
		return
	}
	list, err := a.interactiveServiceClient.BothFollowList(ctx.Request.Context(),
		&interactivev1.BothFollowListRequest{
			Id: info.Id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取互关列表失败",
			Data: "error",
		})
		a.logger.Error("获取互关列表失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:BothFollowList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]int64]{
		Code: 200,
		Msg:  "获取互关列表成功",
		Data: list.Ids,
	})
}

func (a *InteractiveHandler) FanList(ctx *gin.Context) {
	// 获取用户信息
	info, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "InteractiveHandler:FanList"))
		return
	}
	list, err := a.interactiveServiceClient.FanList(ctx.Request.Context(),
		&interactivev1.FanListRequest{
			Id: info.Id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取粉丝列表失败",
			Data: "error",
		})
		a.logger.Error("获取粉丝列表失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveHandler:FanList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]int64]{
		Code: 200,
		Msg:  "获取粉丝列表成功",
		Data: list.Ids,
	})
}
