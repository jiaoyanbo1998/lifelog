package grpc

import (
	"context"
	"github.com/jinzhu/copier"
	collectv1 "lifelog-grpc/api/proto/gen/api/proto/collect/v1"
	lifelogv1 "lifelog-grpc/api/proto/gen/api/proto/lifelog/v1"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/service"
)

type CollectServiceGRPCService struct {
	collectService       service.CollectService
	lifeLogServiceClient lifelogv1.LifeLogServiceClient
	collectv1.UnimplementedCollectServiceServer
}

func NewCollectServiceGRPCService(collectService service.CollectService) *CollectServiceGRPCService {
	return &CollectServiceGRPCService{
		collectService: collectService,
	}
}

func (c *CollectServiceGRPCService) EditCollect(ctx context.Context, request *collectv1.EditCollectRequest) (*collectv1.EditCollectResponse, error) {
	err := c.collectService.EditCollect(ctx, domain.CollectDomain{
		Id:       request.GetCollect().GetCollectId(),
		AuthorId: request.GetCollect().AuthorId,
		Name:     request.GetCollect().GetName(),
	})
	if err != nil {
		return &collectv1.EditCollectResponse{}, err
	}
	return &collectv1.EditCollectResponse{}, nil
}

func (c *CollectServiceGRPCService) DeleteCollect(ctx context.Context, request *collectv1.DeleteCollectRequest) (*collectv1.DeleteCollectResponse, error) {
	err := c.collectService.DeleteCollect(ctx, request.GetIds(), request.GetAuthorId())
	if err != nil {
		return &collectv1.DeleteCollectResponse{}, err
	}
	return &collectv1.DeleteCollectResponse{}, nil
}

func (c *CollectServiceGRPCService) CollectList(ctx context.Context, request *collectv1.CollectListRequest) (*collectv1.CollectListResponse, error) {
	list, err := c.collectService.CollectList(ctx, request.GetAuthorId(), int(request.GetLimit()),
		int(request.GetOffset()))
	if err != nil {
		return &collectv1.CollectListResponse{}, err
	}
	// 将[]domain.CollectDomain，转为[]*CollectDomain
	collects := make([]*collectv1.Collect, 0, len(list))
	err = copier.Copy(&collects, &list)
	if err != nil {
		return &collectv1.CollectListResponse{}, err
	}
	return &collectv1.CollectListResponse{
		Collects: collects,
	}, nil
}

func (c *CollectServiceGRPCService) InsertCollectDetail(ctx context.Context, request *collectv1.InsertCollectDetailRequest) (*collectv1.InsertCollectDetailResponse, error) {
	err := c.collectService.InsertCollectDetail(ctx, domain.CollectDetailDomain{
		CollectId: request.GetCollect().CollectId,
		AuthorId:  request.GetCollect().AuthorId,
		LifeLogId: request.GetCollectDetail().LifeLogId,
	})
	if err != nil {
		return &collectv1.InsertCollectDetailResponse{}, err
	}
	return &collectv1.InsertCollectDetailResponse{}, nil
}

func (c *CollectServiceGRPCService) CollectDetail(ctx context.Context, request *collectv1.CollectDetailRequest) (*collectv1.CollectDetailResponse, error) {
	details, err := c.collectService.CollectDetail(ctx, request.GetCollect().CollectId,
		int(request.GetLimit()),
		int(request.GetOffset()), request.GetCollect().GetAuthorId())
	if err != nil {
		return &collectv1.CollectDetailResponse{}, err
	}
	var lifeLogIds []int64
	for _, val := range details {
		lifeLogIds = append(lifeLogIds, val.LifeLogId)
	}
	res, err := c.lifeLogServiceClient.DetailMany(ctx, &lifelogv1.DetailManyRequest{
		Ids: lifeLogIds,
	})
	if err != nil {
		return &collectv1.CollectDetailResponse{}, err
	}
	// 将[]domain.PublicLifeLogDomain，转为[]*PublicLifeLog
	publicLifeLogs := make([]*collectv1.PublicLifeLog, 0, len(details))
	err = copier.Copy(&publicLifeLogs, &res)
	return &collectv1.CollectDetailResponse{
		CollectDetail: &collectv1.CollectDetail{
			CollectId:  request.GetCollect().CollectId,
			CreateTime: request.GetCollect().CreateTime,
			UpdateTime: request.GetCollect().UpdateTime,
			Status:     request.GetCollect().Status,
			AuthorId:   request.GetCollect().GetStatus(),
		},
		PublicLifeLogs: publicLifeLogs,
	}, nil
}
