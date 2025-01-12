package grpc

import (
	"context"
	"encoding/json"
	commentv1 "lifelog-grpc/api/proto/gen/comment/v1"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/event/feed"
	saramaKafka "lifelog-grpc/comment/event/sarama-kafka"
	"lifelog-grpc/comment/service"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CommentServiceGRPCService struct {
	commentService service.CommentService
	commentv1.UnimplementedCommentServiceServer
	// producer *kafkax.KafkaProducer
	Producer            *saramaKafka.SyncProducer
	logger              loggerx.Logger
	CommentFeedProducer *feed.SyncProducer
}

func NewCommentServiceGRPCService(commentService service.CommentService,
	producer *saramaKafka.SyncProducer, logger loggerx.Logger,
	CommentFeedProducer *feed.SyncProducer) *CommentServiceGRPCService {
	return &CommentServiceGRPCService{
		commentService:      commentService,
		Producer:            producer,
		logger:              logger,
		CommentFeedProducer: CommentFeedProducer,
	}
}

// ProducerCommentEvent sarama的kafka生产者
func (c *CommentServiceGRPCService) ProducerCommentEvent(ctx context.Context, request *commentv1.ProducerCommentEventRequest) (*commentv1.ProducerCommentEventResponse, error) {
	// 调用kafka生产者发送
	err := c.Producer.ProduceCommentEvent(domain.CommentDomain{
		UserId:       request.GetComment().GetUserId(),
		Biz:          request.GetComment().GetBiz(),
		BizId:        request.GetComment().GetBizId(),
		Content:      request.GetComment().GetContent(),
		ParentId:     request.GetComment().GetParentId(),
		RootId:       request.GetComment().GetRootId(),
		TargetUserId: request.GetComment().GetTargetUserId(),
		Uuid:         request.GetComment().Uuid,
	})
	if err != nil {
		c.logger.Error("kafka生产者发送失败", loggerx.Error(err))
		return &commentv1.ProducerCommentEventResponse{}, err
	}
	// 生产评论Feed流
	type commentFeedEvent struct {
		Biz             string `json:"biz"`
		BizId           int64  `json:"biz_id"`
		CommentedUserId int64  `json:"commented_user_id"`
		Content         string `json:"content"`
	}
	ext := commentFeedEvent{
		Biz:             "lifelog",
		BizId:           request.GetComment().GetBizId(),
		CommentedUserId: request.GetComment().GetTargetUserId(),
		Content:         request.GetComment().GetContent(),
	}
	marshal, err := json.Marshal(ext)
	if err != nil {
		c.logger.Error("JSON 序列化失败", loggerx.Error(err))
		return &commentv1.ProducerCommentEventResponse{}, err
	}
	err = c.CommentFeedProducer.ProduceCommentEventFeed(feed.FeedEvent{
		UserId:     request.GetComment().GetUserId(),
		Content:    string(marshal),
		CreateTime: time.Now().UnixMilli(),
		Type:       "LifeLog_comment_event",
	})
	if err != nil {
		c.logger.Error("kafka生产者，feed流失败", loggerx.Error(err))
		return &commentv1.ProducerCommentEventResponse{}, err
	}
	// 返回成功响应
	return &commentv1.ProducerCommentEventResponse{}, nil
}

// ProducerCommentEvent kafka-go的kafka生产者
/*
func (c *CommentServiceGRPCService) ProducerCommentEvent(ctx context.Context, request *commentv1.ProducerCommentEventRequest) (*commentv1.ProducerCommentEventResponse, error) {
	// 定义评论结构体
	type comment struct {
		UserId   int64  `json:"user_id"`
		Biz      string `json:"biz"`
		BizId    int64  `json:"biz_id"`
		Content  string `json:"content"`
		ParentId int64  `json:"parent_id"`
		RootId   int64  `json:"root_id"`
	}
	// 将请求中的评论数据映射到结构体
	com := comment{
		UserId:   request.Comment.GetUserId(),
		Biz:      request.Comment.GetBiz(),
		BizId:    request.Comment.GetBizId(),
		Content:  request.Comment.GetContent(),
		ParentId: request.Comment.GetParentId(),
		RootId:   request.Comment.GetRootId(),
	}
	// json序列化，将数据转为[]byte类型的json对象
	marshal, err := json.Marshal(com)
	if err != nil {
		c.logger.Error("JSON 序列化失败", loggerx.Error(err))
		return nil, status.Errorf(codes.Internal, "JSON 序列化失败: %v", err)
	}
	// 创建Kafka消息
	message := kafka.Message{
		Value: marshal,
	}
	// 发送消息到 Kafka
	c.producer.Send(message)
	// 返回成功响应
	return &commentv1.ProducerCommentEventResponse{}, nil
}
*/

func (c *CommentServiceGRPCService) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := c.commentService.DeleteComment(ctx, request.GetId())
	if err != nil {
		return &commentv1.DeleteCommentResponse{}, err
	}
	return &commentv1.DeleteCommentResponse{}, nil
}

func (c *CommentServiceGRPCService) EditComment(ctx context.Context, request *commentv1.EditCommentRequest) (*commentv1.EditCommentResponse, error) {
	err := c.commentService.EditComment(ctx, domain.CommentDomain{
		Id:      request.GetComment().GetId(),
		Content: request.GetComment().GetContent(),
	})
	if err != nil {
		return &commentv1.EditCommentResponse{}, err
	}
	return &commentv1.EditCommentResponse{}, nil
}

func (c *CommentServiceGRPCService) FirstList(ctx context.Context, request *commentv1.FirstListRequest) (*commentv1.FirstListResponse, error) {
	list, err := c.commentService.FirstList(ctx, request.GetComment().Biz,
		request.GetComment().GetBizId(), request.GetMin())
	if err != nil {
		return &commentv1.FirstListResponse{}, err
	}
	return &commentv1.FirstListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

func (c *CommentServiceGRPCService) EveryRootChildSonList(ctx context.Context, request *commentv1.EveryRootChildSonListRequest) (*commentv1.EveryRootChildSonListResponse, error) {
	list, err := c.commentService.EveryRootChildSonList(ctx, request.GetComment().GetId(),
		request.GetComment().GetRootId(), request.GetLimit())
	if err != nil {
		return &commentv1.EveryRootChildSonListResponse{}, err
	}
	return &commentv1.EveryRootChildSonListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

func (c *CommentServiceGRPCService) SonList(ctx context.Context, request *commentv1.SonListRequest) (*commentv1.SonListResponse, error) {
	list, err := c.commentService.SonList(ctx, request.GetComment().GetParentId(),
		request.GetLimit(), request.GetOffset())
	if err != nil {
		return &commentv1.SonListResponse{}, err
	}
	return &commentv1.SonListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

// 将[]domain.CommonDomain，转为[]*CommentDomain
func (c *CommentServiceGRPCService) toGRPCCommentDomain(lists []domain.CommentDomain) []*commentv1.Comment {
	cs := make([]*commentv1.Comment, 0, len(lists))
	for _, list := range lists {
		cs = append(cs, &commentv1.Comment{
			UserId:     list.UserId,
			Biz:        list.Biz,
			BizId:      list.BizId,
			Content:    list.Content,
			ParentId:   list.ParentId,
			RootId:     list.RootId,
			Id:         list.Id,
			CreateTime: list.CreateTime,
			UpdateTime: list.UpdateTime,
		})
	}
	return cs
}
