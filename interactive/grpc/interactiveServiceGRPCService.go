package grpc

import (
	"context"
	interactivev1 "lifelog-grpc/api/proto/gen/api/proto/interactive/v1"
	"lifelog-grpc/interactive/service"
)

// InteractiveServiceGRPCService 短信服务的grpc服务器
type InteractiveServiceGRPCService struct {
	interactiveService service.InteractiveService
	interactivev1.UnimplementedInteractiveServiceServer
}

func NewCodeServiceGRPCService(interactiveService service.InteractiveService) *InteractiveServiceGRPCService {
	return &InteractiveServiceGRPCService{
		interactiveService: interactiveService,
	}
}

func (i *InteractiveServiceGRPCService) GetInteractiveInfo(ctx context.Context, request *interactivev1.GetInteractiveInfoRequest) (*interactivev1.GetInteractiveInfoResponse, error) {
	info, err := i.GetInteractiveInfo(ctx, &interactivev1.GetInteractiveInfoRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:   request.GetInteractiveDomain().GetBiz(),
			BizId: request.GetInteractiveDomain().GetBizId(),
		},
	})
	if err != nil {
		return &interactivev1.GetInteractiveInfoResponse{}, err
	}
	return &interactivev1.GetInteractiveInfoResponse{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Id:           info.GetInteractiveDomain().GetId(),
			CreateTime:   info.GetInteractiveDomain().GetCreateTime(),
			UpdateTime:   info.GetInteractiveDomain().GetUpdateTime(),
			ReadCount:    info.GetInteractiveDomain().GetReadCount(),
			LikeCount:    info.GetInteractiveDomain().GetLikeCount(),
			CollectCount: info.GetInteractiveDomain().GetCollectCount(),
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
	err := i.interactiveService.IncreaseCollectCount(ctx,
		request.InteractiveDomain.Biz,
		request.GetInteractiveDomain().GetBizId(),
		request.GetInteractiveDomain().GetUserId())
	if err != nil {
		return &interactivev1.CollectResponse{}, err
	}
	return &interactivev1.CollectResponse{}, nil
}

func (i *InteractiveServiceGRPCService) UnCollect(ctx context.Context, request *interactivev1.UnCollectRequest) (*interactivev1.UnCollectResponse, error) {
	err := i.interactiveService.DecreaseCollectCount(ctx,
		request.GetInteractiveDomain().GetBiz(),
		request.GetInteractiveDomain().GetBizId(),
		request.GetInteractiveDomain().GetUserId())
	if err != nil {
		return &interactivev1.UnCollectResponse{}, err
	}
	return &interactivev1.UnCollectResponse{}, nil
}
