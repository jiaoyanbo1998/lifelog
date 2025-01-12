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
	rg.GET("/findCommentFeed", f.FindCommentFeed)
}

func (f *FeedHandler) FindCommentFeed(ctx *gin.Context) {
	type findReq struct {
		UserId int64 `form:"user_id"`
		Limit  int64 `form:"limit"`
	}
	var req findReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		f.logger.Error("参数绑定错误", loggerx.Error(err),
			loggerx.String("method:", "FeedHandler:Find"))
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
			loggerx.String("method:", "FeedHandler:Find"))
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
