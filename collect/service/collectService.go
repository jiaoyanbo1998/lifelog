package service

import (
	"context"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/repository"
)

type CollectService interface {
	EditCollect(ctx context.Context, collectDomain domain.CollectDomain) error
	DeleteCollect(ctx context.Context, ids []int64, authorId int64) error
	CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error)
	InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error
	CollectDetail(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error)
}

type CollectServiceV1 struct {
	collectRepository repository.CollectRepository
}

func NewCollectService(collectRepository repository.CollectRepository) CollectService {
	return &CollectServiceV1{
		collectRepository: collectRepository,
	}
}

// EditCollect 编辑收藏夹
func (c *CollectServiceV1) EditCollect(ctx context.Context, collectDomain domain.CollectDomain) error {
	return c.collectRepository.SaveCollect(ctx, collectDomain)
}

// DeleteCollect 删除收藏夹
func (c *CollectServiceV1) DeleteCollect(ctx context.Context, ids []int64, authorId int64) error {
	return c.collectRepository.DeleteCollect(ctx, ids, authorId)
}

// CollectList 收藏夹列表
func (c *CollectServiceV1) CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error) {
	return c.collectRepository.CollectList(ctx, userId, limit, offset)
}

// InsertCollectDetail 将文章插入收藏夹
func (c *CollectServiceV1) InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error {
	return c.collectRepository.InsertCollectDetail(ctx, detailDomain)
}

// CollectDetail 收藏夹详情
func (c *CollectServiceV1) CollectDetail(ctx context.Context, collectId int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error) {
	return c.collectRepository.CollectDetail(ctx, collectId, limit, offset, authorId)
}
