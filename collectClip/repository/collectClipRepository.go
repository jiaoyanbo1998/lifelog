package repository

import (
	"context"
	"lifelog-grpc/collectClip/domain"
	"lifelog-grpc/collectClip/repository/dao"
)

type CollectClipRepository interface {
	SaveCollectClip(ctx context.Context, collectClipDomain domain.CollectClipDomain) error
	DeleteCollectClip(ctx context.Context, ids []int64) error
	CollectClipList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectClipDomain, error)
	InsertCollectClipDetail(ctx context.Context, detailDomain domain.CollectClipDetailDomain) error
	CollectClipDetail(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectClipDetailDomain, error)
}

type CollectClipRepositoryV1 struct {
	collectClipDao dao.CollectClipDao
}

func NewCollectClipRepository(collectClipDao dao.CollectClipDao) CollectClipRepository {
	return &CollectClipRepositoryV1{
		collectClipDao: collectClipDao,
	}
}

// SaveCollectClip 编辑收藏夹
func (c *CollectClipRepositoryV1) SaveCollectClip(ctx context.Context, collectClipDomain domain.CollectClipDomain) error {
	// 更新收藏夹
	if collectClipDomain.Id > 0 {
		return c.collectClipDao.UpdateCollectClip(ctx, dao.CollectClip{
			Id:     collectClipDomain.Id,
			Name:   collectClipDomain.Name,
			UserId: collectClipDomain.UserId,
		})
	}
	return c.collectClipDao.InsertCollectClip(ctx, dao.CollectClip{
		Name:   collectClipDomain.Name,
		UserId: collectClipDomain.UserId,
	})
}

// DeleteCollectClip 删除收藏夹
func (c *CollectClipRepositoryV1) DeleteCollectClip(ctx context.Context, ids []int64) error {
	return c.collectClipDao.DeleteCollectClipByIds(ctx, ids)
}

// CollectClipList 收藏夹列表
func (c *CollectClipRepositoryV1) CollectClipList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectClipDomain, error) {
	return c.collectClipDao.PageQuery(ctx, userId, limit, offset)
}

// InsertCollectClipDetail 将文章插入到收藏夹
func (c *CollectClipRepositoryV1) InsertCollectClipDetail(ctx context.Context, detailDomain domain.CollectClipDetailDomain) error {
	return c.collectClipDao.InsertCollectClipDetail(ctx, dao.CollectClipDetail{
		LifeLogId: detailDomain.LifeLogId,
		CollectId: detailDomain.CollectId,
	})
}

// CollectClipDetail 收藏夹详情
func (c *CollectClipRepositoryV1) CollectClipDetail(ctx context.Context, id int64, limit int, offset int,
	authorId int64) ([]domain.CollectClipDetailDomain, error) {
	return c.collectClipDao.GetCollectClipDetailById(ctx, id, limit, offset, authorId)
}
