package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	feedv1 "lifelog-grpc/api/proto/gen/feed"
	"lifelog-grpc/errs"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/vo"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"time"
)

type FeedHandler struct {
	feedServiceClient feedv1.FeedServiceClient
	logger            loggerx.Logger
}

func NewFeedHandler(feedServiceClient feedv1.FeedServiceClient, logger loggerx.Logger) *FeedHandler {
	return &FeedHandler{
		feedServiceClient: feedServiceClient,
		logger:            logger,
	}
}

func (f *FeedHandler) RegisterRouters(server *gin.Engine) {
	rg := server.Group("/feed")
	rg.GET("/findLifeLogCommentFeed", f.FindLifeLogCommentFeed)
	rg.GET("/findLifeLogLikeFeed", f.FindLifeLogLikeFeed)
	rg.GET("/findFollowFeed", f.FindFollowFeed)
	rg.GET("/findCollectFeed", f.FindCollectFeed)
	rg.GET("/findReadFeed", f.FindReadFeed)
}

func (f *FeedHandler) FindLifeLogCommentFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"`
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindLifeLogCommentFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	feedEvents, err := f.feedServiceClient.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		UserId:     req.UserId,
		CreateTime: time.Now().UnixMilli(),
		Limit:      req.Limit,
	})
	if err != nil {
		f.logger.Error("查询失败", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindLifeLogCommentFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	fcs := make([]vo.FindCommentFeedVo, 0, len(feedEvents.GetFeedEvents()))
	var lfce domain.LifeLogCommentEvent
	for _, feedEvent := range feedEvents.GetFeedEvents() {
		if feedEvent.GetType() == "LifeLog_comment_event" {
			_ = json.Unmarshal([]byte(feedEvent.Content), &lfce)
			fcs = append(fcs, vo.FindCommentFeedVo{
				UserId:          feedEvent.GetUser().GetId(),
				Biz:             lfce.Biz,
				BizId:           lfce.BizId,
				CommentedUserId: lfce.CommentedUserId,
				Content:         lfce.Content,
			})
		}
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FindCommentFeedVo]{
		Code: 200,
		Msg:  "success",
		Data: fcs,
	})
}

func (f *FeedHandler) FindLifeLogLikeFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"`
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindLifeLogLikeFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	feedEvents, err := f.feedServiceClient.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		UserId:     req.UserId,
		CreateTime: time.Now().UnixMilli(),
		Limit:      req.Limit,
	})
	if err != nil {
		f.logger.Error("查询失败", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindLifeLogLikeFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	fcs := make([]vo.FindLikeFeedVo, 0, len(feedEvents.GetFeedEvents()))
	var lfe domain.LikeFeedEvent
	for _, feedEvent := range feedEvents.GetFeedEvents() {
		if feedEvent.GetType() == "like_event" {
			_ = json.Unmarshal([]byte(feedEvent.Content), &lfe)
			fcs = append(fcs, vo.FindLikeFeedVo{
				UserId:      feedEvent.GetUser().GetId(),
				Biz:         lfe.Biz,
				BizId:       lfe.BizId,
				LikedUserId: lfe.LikedUserId,
			})
		}
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FindLikeFeedVo]{
		Code: 200,
		Msg:  "success",
		Data: fcs,
	})
}

func (f *FeedHandler) FindFollowFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"` // 被关注的人的id
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindFollowFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	// A关注了B，要通知B，A关注了你，去数据库查询id=B_id的feedEvent
	feedEvents, err := f.feedServiceClient.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		UserId:     req.UserId, // 被关注的人的id
		CreateTime: time.Now().UnixMilli(),
		Limit:      req.Limit,
	})
	if err != nil {
		f.logger.Error("查询失败", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindFollowFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	fcs := make([]vo.FindFollowFeedVo, 0, len(feedEvents.GetFeedEvents()))
	var ffe domain.FollowFeedEvent
	for _, feedEvent := range feedEvents.GetFeedEvents() {
		if feedEvent.GetType() == "follow_event" {
			_ = json.Unmarshal([]byte(feedEvent.Content), &ffe)
			fcs = append(fcs, vo.FindFollowFeedVo{
				UserId:         feedEvent.GetUser().GetId(),
				FolloweeUserId: ffe.FolloweeUserId,
				FollowerUserId: ffe.FollowerUserId,
				Biz:            ffe.Biz,
			})
		}
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FindFollowFeedVo]{
		Code: 200,
		Msg:  "success",
		Data: fcs,
	})
}

func (f *FeedHandler) FindCollectFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"` // 被关注的人的id
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindCollectFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	// A收藏了B，要通知B，A收藏了你的文章，去数据库查询id=B_id的feedEvent
	feedEvents, err := f.feedServiceClient.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		UserId:     req.UserId, // 被收藏的lifelog的作者的id
		CreateTime: time.Now().UnixMilli(),
		Limit:      req.Limit,
	})
	if err != nil {
		f.logger.Error("查询失败", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindCollectFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	fcs := make([]vo.FindCollectFeedVo, 0, len(feedEvents.GetFeedEvents()))
	var cfe domain.CollectFeedEvent
	for _, feedEvent := range feedEvents.GetFeedEvents() {
		if feedEvent.GetType() == "collect_event" {
			_ = json.Unmarshal([]byte(feedEvent.Content), &cfe)
			fcs = append(fcs, vo.FindCollectFeedVo{
				UserId:          cfe.UserId, // 收藏的lifeLog的作者的id
				Biz:             cfe.Biz,
				BizId:           cfe.BizId,
				CollectedUserId: cfe.CollectedUserId,
			})
		}
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FindCollectFeedVo]{
		Code: 200,
		Msg:  "success",
		Data: fcs,
	})
}

func (f *FeedHandler) FindReadFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"` // 被关注的人的id
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindReadFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	// A阅读了B的lifelog，要通知B，A阅读了你的lifelog，去数据库查询id=B_id的feedEvent
	feedEvents, err := f.feedServiceClient.FindFeedEvents(ctx, &feedv1.FindFeedEventsRequest{
		UserId:     req.UserId, // 被阅读的lifelog的作者的id
		CreateTime: time.Now().UnixMilli(),
		Limit:      req.Limit,
	})
	if err != nil {
		f.logger.Error("查询失败", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:FindReadFeed"))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	fcs := make([]vo.FindReadFeedVo, 0, len(feedEvents.GetFeedEvents()))
	var rfe domain.ReadFeedEvent
	for _, feedEvent := range feedEvents.GetFeedEvents() {
		if feedEvent.GetType() == "read_event" {
			_ = json.Unmarshal([]byte(feedEvent.Content), &rfe)
			fcs = append(fcs, vo.FindReadFeedVo{
				UserId:       rfe.UserId,       // 粉丝id
				Biz:          rfe.Biz,          // lifelog
				BizId:        rfe.BizId,        // lifelog的id
				ReadedUserId: rfe.ReadedUserId, // 作者的id
			})
		}
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FindReadFeedVo]{
		Code: 200,
		Msg:  "success",
		Data: fcs,
	})
}
