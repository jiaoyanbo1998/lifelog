package repository

import (
	"context"
	"lifelog-grpc/interactive/domain"
	"lifelog-grpc/interactive/repository/cache"
	"lifelog-grpc/interactive/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64, userId int64) error
	IncreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	IncreaseCollectCount(ctx context.Context, biz string, bizId int64, userId, collectId int64) error
	DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId, collectId int64) error
	BatchInteractiveReadCount(ctx context.Context, biz string, bizIds, userIds []int64) error
	GetInteractiveInfoByBizId(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error)
	InsertFollow(ctx context.Context, followerId, followeeId int64) error
	CancelFollow(ctx context.Context, followerId, followeeId int64) error
	FollowList(ctx context.Context, id int64) ([]int64, error)
	FanList(ctx context.Context, id int64) ([]int64, error)
	BothFollowList(ctx context.Context, id int64) ([]int64, error)
}

type InteractiveRepositoryV1 struct {
	interactiveDao   dao.InteractiveDao
	interactiveCache cache.InteractiveCache
}

func NewInteractiveRepository(interactiveDao dao.InteractiveDao,
	interactiveCache cache.InteractiveCache) InteractiveRepository {
	return &InteractiveRepositoryV1{
		interactiveDao:   interactiveDao,
		interactiveCache: interactiveCache,
	}
}

func (i *InteractiveRepositoryV1) FollowList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveDao.FollowList(ctx, id)
}

func (i *InteractiveRepositoryV1) FanList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveDao.FanList(ctx, id)
}

func (i *InteractiveRepositoryV1) BothFollowList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveDao.BothFollowList(ctx, id)
}

func (i *InteractiveRepositoryV1) InsertFollow(ctx context.Context, followerId, followeeId int64) error {
	return i.interactiveDao.InsertFollow(ctx, followerId, followeeId)
}

func (i *InteractiveRepositoryV1) CancelFollow(ctx context.Context, followerId, followeeId int64) error {
	return i.interactiveDao.CancelFollow(ctx, followerId, followeeId)
}

// IncreaseReadCount 增加阅读数
func (i *InteractiveRepositoryV1) IncreaseReadCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 先更新数据库中的阅读数
	err := i.interactiveDao.InsertReadCount(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	// 再删除缓存中的阅读数
	err = i.interactiveCache.DeleteReadCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

// IncreaseLikeCount 增加点赞数
func (i *InteractiveRepositoryV1) IncreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 先更新数据库中的点赞数
	err := i.interactiveDao.InsertLikeCount(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	// 再更新缓存中的点赞数
	err = i.interactiveCache.InsertLikeCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseLikeCount 减少点赞数
func (i *InteractiveRepositoryV1) DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 先更新数据库中的点赞数
	err := i.interactiveDao.DecreaseLikeCount(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	// 再更新缓存中的点赞数
	err = i.interactiveCache.DecreaseLikeCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

// IncreaseCollectCount 增加收藏数
func (i *InteractiveRepositoryV1) IncreaseCollectCount(ctx context.Context, biz string,
	bizId int64, userId, collectId int64) error {
	// 先更新数据库中的收藏数
	err := i.interactiveDao.InsertCollectCount(ctx, biz, bizId, userId, collectId)
	if err != nil {
		return err
	}
	// 再更新缓存中的收藏数
	err = i.interactiveCache.InsertCollectCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseCollectCount 减少收藏数
func (i *InteractiveRepositoryV1) DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId, collectId int64) error {
	// 先操作数据库
	err := i.interactiveDao.DecreaseCollectCount(ctx, biz, bizId, userId, collectId)
	if err != nil {
		return err
	}
	// 再操作redis缓存
	err = i.interactiveCache.DecreaseCollectCount(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

// BatchInteractiveReadCount 批量增加阅读数
func (i *InteractiveRepositoryV1) BatchInteractiveReadCount(ctx context.Context,
	biz string, bizIds, userIds []int64) error {
	return i.interactiveDao.BatchInteractiveReadCount(ctx, biz, bizIds, userIds)
}

// GetInteractiveInfoByBizId 根据文章id获取文章的互动信息
func (i *InteractiveRepositoryV1) GetInteractiveInfoByBizId(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error) {
	return i.interactiveDao.GetInteractiveInfoByBizId(ctx, biz, bizId)
}
