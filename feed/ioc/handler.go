package ioc

import (
	"lifelog-grpc/feed/repository"
	"lifelog-grpc/feed/service"
)

// RegisterHandler 注册handler
func RegisterHandler(feedRepository repository.FeedRepository) map[string]service.Handler {
	followHandler := service.NewFollowEventHandler(feedRepository)
	likeHandler := service.NewLikeEventHandler(feedRepository)
	commentHandler := service.NewCommentEventHandler(feedRepository)
	collectHandler := service.NewCollectEventHandler(feedRepository)
	readHandler := service.NewReadEventHandler(feedRepository)
	return map[string]service.Handler{
		service.FollowEventName:  followHandler,
		service.LikeEventName:    likeHandler,
		service.CommentEventName: commentHandler,
		service.CollectEventName: collectHandler,
		service.ReadEventName:    readHandler,
	}
}
