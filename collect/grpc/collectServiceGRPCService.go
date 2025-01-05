package grpc

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	collectv1 "lifelog-grpc/api/proto/gen/collect/v1"
	lifelogv1 "lifelog-grpc/api/proto/gen/lifelog/v1"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/service"
)

type CollectServiceGRPCService struct {
	collectService       service.CollectService
	lifeLogServiceClient lifelogv1.LifeLogServiceClient
	collectv1.UnimplementedCollectServiceServer
}

func NewCollectServiceGRPCService(collectService service.CollectService,
	lifeLogServiceClient lifelogv1.LifeLogServiceClient) *CollectServiceGRPCService {
	return &CollectServiceGRPCService{
		collectService:       collectService,
		lifeLogServiceClient: lifeLogServiceClient,
	}
}

func (c *CollectServiceGRPCService) DeleteCollectDetail(ctx context.Context, request *collectv1.DeleteCollectDetailRequest) (*collectv1.DeleteCollectDetailResponse, error) {
	err := c.collectService.DeleteCollectDetail(ctx,
		request.CollectId, request.LifeLogId, request.AuthorId)
	if err != nil {
		return &collectv1.DeleteCollectDetailResponse{}, err
	}
	return &collectv1.DeleteCollectDetailResponse{}, nil
}

func (c *CollectServiceGRPCService) EditCollect(ctx context.Context, request *collectv1.EditCollectRequest) (*collectv1.EditCollectResponse, error) {
	err := c.collectService.EditCollect(ctx, domain.CollectDomain{
		Id:       request.GetCollect().GetId(),
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
		CollectId: request.GetCollect().GetId(),
		AuthorId:  request.GetCollect().AuthorId,
		LifeLogId: request.GetCollectDetail().LifeLogId,
	})
	if err != nil {
		return &collectv1.InsertCollectDetailResponse{}, err
	}
	return &collectv1.InsertCollectDetailResponse{}, nil
}

func (c *CollectServiceGRPCService) CollectDetail(ctx context.Context, request *collectv1.CollectDetailRequest) (*collectv1.CollectDetailResponse, error) {
	details, err := c.collectService.CollectDetail(ctx,
		request.GetCollect().GetId(),
		int(request.GetLimit()),
		int(request.GetOffset()), request.GetCollect().GetAuthorId())
	if err != nil {
		return &collectv1.CollectDetailResponse{}, err
	}
	if len(details) == 0 {
		return &collectv1.CollectDetailResponse{}, fmt.Errorf("no details found")
	}
	var lifeLogIds []int64
	for _, val := range details {
		lifeLogIds = append(lifeLogIds, val.LifeLogId)
	}
	// 如果生命日志ID为空，直接返回
	if len(lifeLogIds) == 0 {
		return &collectv1.CollectDetailResponse{
			CollectDetail: &collectv1.CollectDetail{
				CollectId:  details[0].CollectId,
				CreateTime: details[0].CreateTime,
				UpdateTime: details[0].UpdateTime,
				Status:     int64(details[0].Status),
				AuthorId:   details[0].AuthorId,
			},
			PublicLifeLogs: []*collectv1.PublicLifeLog{},
		}, nil
	}
	res, err := c.lifeLogServiceClient.DetailMany(ctx, &lifelogv1.DetailManyRequest{
		Ids: lifeLogIds,
	})
	if err != nil {
		return &collectv1.CollectDetailResponse{}, err
	}
	// 将[]domain.PublicLifeLogDomain，转为[]*PublicLifeLog
	publicLifeLogs := make([]*collectv1.PublicLifeLog, 0, len(res.GetLifeLogDomain()))
	for _, val := range res.GetLifeLogDomain() {
		publicLifeLogs = append(publicLifeLogs, &collectv1.PublicLifeLog{
			PublicLifeLogId: val.Id,
			Title:           val.GetTitle(),
			Content:         val.GetContent(),
			AuthorId:        val.Author.GetUserId(),
			CreateTime:      val.GetCreateTime(),
			UpdateTime:      val.GetUpdateTime(),
			Status:          val.GetStatus(),
		})
	}
	return &collectv1.CollectDetailResponse{
		CollectDetail: &collectv1.CollectDetail{
			CollectId:  details[0].CollectId,
			CreateTime: details[0].CreateTime,
			UpdateTime: details[0].UpdateTime,
			Status:     int64(details[0].Status),
			AuthorId:   details[0].AuthorId,
		},
		PublicLifeLogs: publicLifeLogs,
	}, nil
}
