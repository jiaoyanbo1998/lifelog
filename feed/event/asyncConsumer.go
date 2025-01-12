package event

import (
	"context"
	"github.com/IBM/sarama"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
	"time"
)

type FeedEventAsyncConsumer struct {
	client      sarama.Client
	logger      loggerx.Logger
	feedService service.FeedService
}

func NewFeedEventAsyncConsumer(
	client sarama.Client,
	logger loggerx.Logger,
	feedService service.FeedService) *FeedEventAsyncConsumer {
	return &FeedEventAsyncConsumer{
		feedService: feedService,
		client:      client,
		logger:      logger,
	}
}

func (r *FeedEventAsyncConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("feed_event",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{"LifeLog_comment_event_feed"}, // 消费的topic
			saramax.NewHandler[domain.FeedEvent](r.logger, r.Consume))
		if er != nil {
			r.logger.Error("退出了消费循环异常", loggerx.Error(err))
		}
	}()
	return err
}

func (r *FeedEventAsyncConsumer) Consume(msg *sarama.ConsumerMessage, feedEvent domain.FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	de := domain.FeedEvent{
		Type:       feedEvent.Type,
		UserId:     feedEvent.UserId,
		Content:    feedEvent.Content,
		CreateTime: feedEvent.CreateTime,
	}
	return r.feedService.CreateFeedEvent(ctx, de)
}
