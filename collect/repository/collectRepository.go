package repository

import (
	"context"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/repository/dao"
)

type CollectRepository interface {
	SaveCollect(ctx context.Context, collectDomain domain.CollectDomain) error
	DeleteCollect(ctx context.Context, ids []int64) error
	CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error)
	InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error
	CollectDetail(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error)
}

type CollectRepositoryV1 struct {
	collectDao dao.CollectDao
}

func NewCollectRepository(collectDao dao.CollectDao) CollectRepository {
	return &CollectRepositoryV1{
		collectDao: collectDao,
	}
}

// SaveCollect 编辑收藏夹
func (c *CollectRepositoryV1) SaveCollect(ctx context.Context, collectDomain domain.CollectDomain) error {
	// 更新收藏夹
	if collectDomain.Id > 0 {
		return c.collectDao.UpdateCollect(ctx, dao.Collect{
			Id:     collectDomain.Id,
			Name:   collectDomain.Name,
			UserId: collectDomain.UserId,
		})
	}
	return c.collectDao.InsertCollect(ctx, dao.Collect{
		Name:   collectDomain.Name,
		UserId: collectDomain.UserId,
	})
}

// DeleteCollect 删除收藏夹
func (c *CollectRepositoryV1) DeleteCollect(ctx context.Context, ids []int64) error {
	return c.collectDao.DeleteCollectByIds(ctx, ids)
}

// CollectList 收藏夹列表
func (c *CollectRepositoryV1) CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error) {
	return c.collectDao.PageQuery(ctx, userId, limit, offset)
}

// InsertCollectDetail 将文章插入到收藏夹
func (c *CollectRepositoryV1) InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error {
	return c.collectDao.InsertCollectDetail(ctx, dao.CollectDetail{
		LifeLogId: detailDomain.LifeLogId,
		CollectId: detailDomain.CollectId,
	})
}

// CollectDetail 收藏夹详情
func (c *CollectRepositoryV1) CollectDetail(ctx context.Context, id int64, limit int, offset int,
	authorId int64) ([]domain.CollectDetailDomain, error) {
	return c.collectDao.GetCollectDetailById(ctx, id, limit, offset, authorId)
}
