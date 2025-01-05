package grpc

import (
	"context"
	"github.com/jinzhu/copier"
	lifelogv1 "lifelog-grpc/api/proto/gen/api/proto/lifelog/v1"
	"lifelog-grpc/lifeLog/domain"
	"lifelog-grpc/lifeLog/service"
	"lifelog-grpc/pkg/loggerx"
)

type LifeLogServiceGRPCService struct {
	lifeLogService service.LifeLogService
	logger         loggerx.Logger
	lifelogv1.UnimplementedLifeLogServiceServer
}

func NewLifeLogServiceGRPCService(lifeLogService service.LifeLogService,
	logger loggerx.Logger) *LifeLogServiceGRPCService {
	return &LifeLogServiceGRPCService{
		lifeLogService: lifeLogService,
		logger:         logger,
	}
}

func (l *LifeLogServiceGRPCService) DetailMany(ctx context.Context, request *lifelogv1.DetailManyRequest) (*lifelogv1.DetailManyResponse, error) {
	res, err := l.lifeLogService.DetailMany(ctx, request.GetIds())
	if err != nil {
		return &lifelogv1.DetailManyResponse{}, err
	}
	// 将[]domain.LifeLogDomain，转为[]*LifeLogDomain
	llds := make([]*lifelogv1.LifeLogDomain, 0, len(res))
	for _, v := range res {
		llds = append(llds, &lifelogv1.LifeLogDomain{
			Id:         v.Id,
			Title:      v.Title,
			Content:    v.Content,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
			Status:     int64(v.Status),
			Author: &lifelogv1.Author{
				UserId: v.Author.Id,
			},
		})
	}
	if err != nil {
		l.logger.Error("copier失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogServiceGRPCService:DetailMany"))
		return &lifelogv1.DetailManyResponse{}, err
	}
	return &lifelogv1.DetailManyResponse{
		LifeLogDomain: llds,
	}, nil
}

func (l *LifeLogServiceGRPCService) Detail(ctx context.Context, request *lifelogv1.DetailRequest) (*lifelogv1.DetailResponse, error) {
	detail, err := l.lifeLogService.Detail(ctx, request.GetLifeLogDomain().GetId(), request.GetIsPublic())
	if err != nil {
		return &lifelogv1.DetailResponse{}, err
	}
	return &lifelogv1.DetailResponse{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id:         detail.Id,
			Title:      detail.Title,
			Content:    detail.Content,
			CreateTime: detail.CreateTime,
			UpdateTime: detail.UpdateTime,
			Status:     int64(detail.Status),
			Author: &lifelogv1.Author{
				UserId:   detail.Author.Id,
				NickName: detail.Author.Name,
			},
		},
	}, nil
}

func (l *LifeLogServiceGRPCService) Edit(ctx context.Context, request *lifelogv1.EditRequest) (*lifelogv1.EditResponse, error) {
	// 调用service层
	res, err := l.lifeLogService.Save(ctx, domain.LifeLogDomain{
		Id:      request.GetLifeLogDomain().GetId(),
		Title:   request.GetLifeLogDomain().GetTitle(),
		Content: request.GetLifeLogDomain().GetContent(),
		Author: domain.Author{
			Id: request.GetLifeLogDomain().GetAuthor().GetUserId(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &lifelogv1.EditResponse{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id:         res.Id,
			Title:      res.Title,
			Content:    res.Content,
			CreateTime: res.CreateTime,
			UpdateTime: res.UpdateTime,
			Author: &lifelogv1.Author{
				UserId:   res.Author.Id,
				NickName: res.Author.Name,
			},
		},
	}, nil
}

func (l *LifeLogServiceGRPCService) Delete(ctx context.Context, request *lifelogv1.DeleteRequest) (*lifelogv1.DeleteResponse, error) {
	err := l.lifeLogService.Delete(ctx, request.GetIds(), request.GetIsPublic())
	if err != nil {
		return &lifelogv1.DeleteResponse{}, err
	}
	return &lifelogv1.DeleteResponse{}, nil
}

func (l *LifeLogServiceGRPCService) SearchByTitle(ctx context.Context, request *lifelogv1.SearchByTitleRequest) (*lifelogv1.SearchByTitleResponse, error) {
	res, err := l.lifeLogService.SearchByTitle(ctx, request.GetLifeLogDomain().GetTitle(),
		request.GetLifeLogDomain().GetLimit(),
		request.GetLifeLogDomain().GetOffset())
	if err != nil {
		return &lifelogv1.SearchByTitleResponse{}, err
	}
	// 将[]domain.LifeLogDomain，转为[]*LifeLogDomain
	llds := make([]*lifelogv1.LifeLogDomain, 0, len(res))
	err = copier.Copy(&llds, &res)
	if err != nil {
		l.logger.Error("copier失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogServiceGRPCService:SearchByTitle"))
		return &lifelogv1.SearchByTitleResponse{}, err
	}
	return &lifelogv1.SearchByTitleResponse{
		LifeLogDomain: llds,
	}, nil
}

func (l *LifeLogServiceGRPCService) DraftList(ctx context.Context, request *lifelogv1.DraftListRequest) (*lifelogv1.DraftListResponse, error) {
	res, err := l.lifeLogService.SearchByAuthorId(ctx,
		request.GetLifeLogDomain().GetAuthor().GetUserId(),
		request.GetLifeLogDomain().GetLimit(),
		request.GetLifeLogDomain().GetOffset())
	if err != nil {
		return &lifelogv1.DraftListResponse{}, err
	}
	// 将[]domain.LifeLogDomain，转为[]*LifeLogDomain
	llds := make([]*lifelogv1.LifeLogDomain, 0, len(res))
	err = copier.Copy(&llds, &res)
	if err != nil {
		l.logger.Error("copier失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogServiceGRPCService:SearchByTitle"))
		return &lifelogv1.DraftListResponse{}, err
	}
	return &lifelogv1.DraftListResponse{
		LifeLogDomain: llds,
	}, nil
}

func (l *LifeLogServiceGRPCService) Revoke(ctx context.Context, request *lifelogv1.RevokeRequest) (*lifelogv1.RevokeResponse, error) {
	err := l.lifeLogService.Revoke(ctx, request.GetLifeLogDomain().GetId(),
		request.GetLifeLogDomain().GetAuthor().GetUserId())
	if err != nil {
		return &lifelogv1.RevokeResponse{}, err
	}
	return &lifelogv1.RevokeResponse{}, nil
}

func (l *LifeLogServiceGRPCService) Publish(ctx context.Context, request *lifelogv1.PublishRequest) (*lifelogv1.PublishResponse, error) {
	err := l.lifeLogService.Publish(ctx, domain.LifeLogDomain{
		Id:      request.GetLifeLogDomain().GetId(),
		Title:   request.GetLifeLogDomain().GetTitle(),
		Content: request.GetLifeLogDomain().GetContent(),
		Author: domain.Author{
			Id: request.GetLifeLogDomain().GetAuthor().GetUserId(),
		},
	})
	if err != nil {
		return &lifelogv1.PublishResponse{}, err
	}
	return &lifelogv1.PublishResponse{}, nil
}
