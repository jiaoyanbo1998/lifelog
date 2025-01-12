package grpc

import (
	"context"
	feedv1 "lifelog-grpc/api/proto/gen/feed"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/service"
)

type FeedServiceGRPCService struct {
	feedService service.FeedService
	feedv1.UnimplementedFeedServiceServer
}

func NewFeedServiceGRPCService(feedService service.FeedService) *FeedServiceGRPCService {
	return &FeedServiceGRPCService{feedService: feedService}
}

func (f *FeedServiceGRPCService) FindFeedEvents(ctx context.Context, request *feedv1.FindFeedEventsRequest) (*feedv1.FindFeedEventsResponse, error) {
	findFeedEvents, err := f.feedService.FindFeedEvents(ctx, request.GetUserId(),
		request.GetCreateTime(), request.GetLimit())
	if err != nil {
		return &feedv1.FindFeedEventsResponse{}, err
	}
	// 将[]domain.FeedEvent转换为[]*feedv1.FeedEvent
	fds := make([]*feedv1.FeedEvent, 0, len(findFeedEvents))
	for _, fd := range findFeedEvents {
		fds = append(fds, f.convertToGRPC(fd))
	}
	return &feedv1.FindFeedEventsResponse{
		FeedEvents: fds,
	}, nil
}

func (f *FeedServiceGRPCService) convertToGRPC(feedEvent domain.FeedEvent) *feedv1.FeedEvent {
	return &feedv1.FeedEvent{
		Id:         feedEvent.ID,
		Type:       feedEvent.Type,
		CreateTime: feedEvent.CreateTime,
		Content:    feedEvent.Content,
		User: &feedv1.User{
			Id: feedEvent.UserId,
		},
	}
}
