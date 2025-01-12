package grpc

import (
	"context"
	"encoding/json"
	collectv1 "lifelog-grpc/api/proto/gen/collect/v1"
	feedv1 "lifelog-grpc/api/proto/gen/feed"
	interactivev1 "lifelog-grpc/api/proto/gen/interactive/v1"
	"lifelog-grpc/interactive/event/feed"
	"lifelog-grpc/interactive/service"
	"time"
)

// InteractiveServiceGRPCService 短信服务的grpc服务器
type InteractiveServiceGRPCService struct {
	interactiveService   service.InteractiveService
	collectServiceClient collectv1.CollectServiceClient
	interactivev1.UnimplementedInteractiveServiceServer
	feedServiceClient   feedv1.FeedServiceClient
	interactiveProducer *feed.SyncProducer
}

func NewCodeServiceGRPCService(interactiveService service.InteractiveService,
	collectServiceClient collectv1.CollectServiceClient,
	feedServiceClient feedv1.FeedServiceClient,
	interactiveProducer *feed.SyncProducer) *InteractiveServiceGRPCService {
	return &InteractiveServiceGRPCService{
		interactiveService:   interactiveService,
		collectServiceClient: collectServiceClient,
		feedServiceClient:    feedServiceClient,
		interactiveProducer:  interactiveProducer,
	}
}

func (i *InteractiveServiceGRPCService) FollowList(ctx context.Context, request *interactivev1.FollowListRequest) (*interactivev1.FollowListResponse, error) {
	list, err := i.interactiveService.FollowList(ctx, request.Id)
	if err != nil {
		return &interactivev1.FollowListResponse{}, err
	}
	return &interactivev1.FollowListResponse{
		Ids: list,
	}, nil
}

func (i *InteractiveServiceGRPCService) FanList(ctx context.Context, request *interactivev1.FanListRequest) (*interactivev1.FanListResponse, error) {
	list, err := i.interactiveService.FanList(ctx, request.Id)
	if err != nil {
		return &interactivev1.FanListResponse{}, err
	}
	return &interactivev1.FanListResponse{
		Ids: list,
	}, nil
}

func (i *InteractiveServiceGRPCService) BothFollowList(ctx context.Context, request *interactivev1.BothFollowListRequest) (*interactivev1.BothFollowListResponse, error) {
	list, err := i.interactiveService.BothFollowList(ctx, request.Id)
	if err != nil {
		return &interactivev1.BothFollowListResponse{}, err
	}
	return &interactivev1.BothFollowListResponse{
		Ids: list,
	}, nil
}

func (i *InteractiveServiceGRPCService) InsertFollow(ctx context.Context, request *interactivev1.InsertFollowRequest) (*interactivev1.InsertFollowResponse, error) {
	err := i.interactiveService.InsertFollow(ctx, request.GetFollow().GetFollowerId(),
		request.GetFollow().GetFolloweeId())
	if err != nil {
		return &interactivev1.InsertFollowResponse{}, err
	}
	return &interactivev1.InsertFollowResponse{}, nil
}

func (i *InteractiveServiceGRPCService) CancelFollow(ctx context.Context, request *interactivev1.CancelFollowRequest) (*interactivev1.CancelFollowResponse, error) {
	err := i.interactiveService.CancelFollow(ctx, request.GetFollow().GetFollowerId(),
		request.GetFollow().GetFolloweeId())
	if err != nil {
		return &interactivev1.CancelFollowResponse{}, err
	}
	return &interactivev1.CancelFollowResponse{}, nil
}

func (i *InteractiveServiceGRPCService) IncreaseRead(ctx context.Context, request *interactivev1.IncreaseReadRequest) (*interactivev1.IncreaseReadResponse, error) {
	err := i.interactiveService.IncreaseReadCount(ctx,
		request.GetInteractiveDomain().GetBiz(),
		request.GetInteractiveDomain().GetBizId(),
		request.GetInteractiveDomain().GetUserId())
	if err != nil {
		return &interactivev1.IncreaseReadResponse{}, err
	}
	return &interactivev1.IncreaseReadResponse{}, nil
}

func (i *InteractiveServiceGRPCService) InteractiveInfo(
	ctx context.Context, request *interactivev1.InteractiveInfoRequest) (*interactivev1.InteractiveInfoResponse, error) {
	info, err := i.interactiveService.GetInteractiveInfo(ctx,
		request.GetInteractiveDomain().GetBiz(),
		request.GetInteractiveDomain().GetBizId(),
	)
	if err != nil {
		return &interactivev1.InteractiveInfoResponse{}, err
	}
	return &interactivev1.InteractiveInfoResponse{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Id:           info.Id,
			CreateTime:   info.CreateTime,
			UpdateTime:   info.UpdateTime,
			ReadCount:    info.ReadCount,
			LikeCount:    info.LikeCount,
			CollectCount: info.CollectCount,
		},
	}, nil
}

func (i *InteractiveServiceGRPCService) Like(ctx context.Context, request *interactivev1.LikeRequest) (*interactivev1.LikeResponse, error) {
	err := i.interactiveService.IncreaseLikeCount(
		ctx, request.InteractiveDomain.Biz,
		request.InteractiveDomain.BizId,
		request.InteractiveDomain.UserId)
	if err != nil {
		return &interactivev1.LikeResponse{}, err
	}
	// 生产评论Feed流
	type likeFeedEvent struct {
		Biz         string `json:"biz"`
		BizId       int64  `json:"biz_id"`
		LikedUserId int64  `json:"liked_user_id"`
	}
	ext := likeFeedEvent{
		Biz:         "lifelog",
		BizId:       request.GetInteractiveDomain().GetBizId(),
		LikedUserId: request.GetInteractiveDomain().GetTargetUserId(),
	}
	marshal, err := json.Marshal(ext)
	if err != nil {
		return &interactivev1.LikeResponse{}, err
	}
	err = i.interactiveProducer.ProduceInteractiveEventFeed(feed.FeedEvent{
		UserId:     request.GetInteractiveDomain().GetUserId(),
		Content:    string(marshal),
		CreateTime: time.Now().UnixMilli(),
		Type:       "like_event",
	})
	if err != nil {
		return &interactivev1.LikeResponse{}, err
	}
	return &interactivev1.LikeResponse{}, nil
}

func (i *InteractiveServiceGRPCService) UnLike(ctx context.Context, request *interactivev1.UnLikeRequest) (*interactivev1.UnLikeResponse, error) {
	err := i.interactiveService.DecreaseLikeCount(ctx, request.InteractiveDomain.Biz,
		request.InteractiveDomain.BizId,
		request.InteractiveDomain.UserId)
	if err != nil {
		return &interactivev1.UnLikeResponse{}, err
	}
	return &interactivev1.UnLikeResponse{}, nil
}

func (i *InteractiveServiceGRPCService) Collect(ctx context.Context, request *interactivev1.CollectRequest) (*interactivev1.CollectResponse, error) {
	// 操作互动表
	err := i.interactiveService.IncreaseCollectCount(ctx,
		request.InteractiveDomain.Biz,
		request.GetInteractiveDomain().GetBizId(),
		request.GetInteractiveDomain().GetUserId(),
		request.GetCollectId())
	if err != nil {
		return &interactivev1.CollectResponse{}, err
	}
	// 调用收藏夹服务
	// 插入收藏详情表
	_, err = i.collectServiceClient.InsertCollectDetail(ctx, &collectv1.InsertCollectDetailRequest{
		Collect: &collectv1.Collect{
			Id:       request.GetCollectId(),
			AuthorId: request.GetInteractiveDomain().GetUserId(),
		},
		CollectDetail: &collectv1.CollectDetail{
			LifeLogId: request.GetInteractiveDomain().GetBizId(),
		},
	})
	if err != nil {
		return &interactivev1.CollectResponse{}, err
	}
	return &interactivev1.CollectResponse{}, nil
}

func (i *InteractiveServiceGRPCService) UnCollect(ctx context.Context, request *interactivev1.UnCollectRequest) (*interactivev1.UnCollectResponse, error) {
	err := i.interactiveService.DecreaseCollectCount(ctx,
		request.GetInteractiveDomain().GetBiz(),
		request.GetInteractiveDomain().GetBizId(),
		request.GetInteractiveDomain().GetUserId(),
		request.GetCollectId())
	if err != nil {
		return &interactivev1.UnCollectResponse{}, err
	}
	// 调用收藏夹服务
	// 取消收藏后，要将这个lifelog从收藏夹中详情表中移除此条记录
	_, err = i.collectServiceClient.DeleteCollectDetail(ctx, &collectv1.DeleteCollectDetailRequest{
		CollectId: request.CollectId,
		LifeLogId: request.GetInteractiveDomain().GetBizId(),
		AuthorId:  request.GetInteractiveDomain().GetUserId(),
	})
	if err != nil {
		return &interactivev1.UnCollectResponse{}, err
	}
	return &interactivev1.UnCollectResponse{}, nil
}
