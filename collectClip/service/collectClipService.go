package service

import (
	"context"
	"lifelog-grpc/collectClip/domain"
	"lifelog-grpc/collectClip/repository"
)

type CollectClipService interface {
	EditCollectClip(ctx context.Context, collectClipDomain domain.CollectClipDomain) error
	DeleteCollectClip(ctx context.Context, ids []int64) error
	CollectClipList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectClipDomain, error)
	InsertCollectClipDetail(ctx context.Context, detailDomain domain.CollectClipDetailDomain) error
	CollectClipDetail(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectClipDetailDomain, error)
}

type CollectClipServiceV1 struct {
	collectClipRepository repository.CollectClipRepository
}

func NewCollectClipService(collectClipRepository repository.CollectClipRepository) CollectClipService {
	return &CollectClipServiceV1{
		collectClipRepository: collectClipRepository,
	}
}

// EditCollectClip 编辑收藏夹
func (c *CollectClipServiceV1) EditCollectClip(ctx context.Context, collectClipDomain domain.CollectClipDomain) error {
	return c.collectClipRepository.SaveCollectClip(ctx, collectClipDomain)
}

// DeleteCollectClip 删除收藏夹
func (c *CollectClipServiceV1) DeleteCollectClip(ctx context.Context, ids []int64) error {
	return c.collectClipRepository.DeleteCollectClip(ctx, ids)
}

// CollectClipList 收藏夹列表
func (c *CollectClipServiceV1) CollectClipList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectClipDomain, error) {
	return c.collectClipRepository.CollectClipList(ctx, userId, limit, offset)
}

// InsertCollectClipDetail 将文章插入收藏夹
func (c *CollectClipServiceV1) InsertCollectClipDetail(ctx context.Context, detailDomain domain.CollectClipDetailDomain) error {
	return c.collectClipRepository.InsertCollectClipDetail(ctx, detailDomain)
}

// CollectClipDetail 收藏夹详情
func (c *CollectClipServiceV1) CollectClipDetail(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectClipDetailDomain, error) {
	return c.collectClipRepository.CollectClipDetail(ctx, id, limit, offset, authorId)
}
