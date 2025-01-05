package service

import (
	"context"
	"lifelog-grpc/interactive/domain"
	"lifelog-grpc/interactive/repository"
)

type InteractiveService interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64, userId int64) error
	BatchInteractiveReadCount(ctx context.Context, biz string, bizIds, userIds []int64) error
	IncreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	IncreaseCollectCount(ctx context.Context, biz string, bizId int64, userId int64, collectId int64) error
	DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId int64, collectId int64) error
	GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error)
	InsertFollow(ctx context.Context, followerId, followeeId int64) error
	CancelFollow(ctx context.Context, followerId, followeeId int64) error
	FollowList(ctx context.Context, id int64) ([]int64, error)
	FanList(ctx context.Context, id int64) ([]int64, error)
	BothFollowList(ctx context.Context, id int64) ([]int64, error)
}

type InteractiveServiceV1 struct {
	interactiveRepository repository.InteractiveRepository
}

func NewInteractiveService(interactiveRepository repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceV1{
		interactiveRepository: interactiveRepository,
	}
}

func (i *InteractiveServiceV1) FollowList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveRepository.FollowList(ctx, id)
}

func (i *InteractiveServiceV1) FanList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveRepository.FanList(ctx, id)
}

func (i *InteractiveServiceV1) BothFollowList(ctx context.Context, id int64) ([]int64, error) {
	return i.interactiveRepository.BothFollowList(ctx, id)
}

func (i *InteractiveServiceV1) InsertFollow(ctx context.Context, followerId, followeeId int64) error {
	return i.interactiveRepository.InsertFollow(ctx, followerId, followeeId)
}

func (i *InteractiveServiceV1) CancelFollow(ctx context.Context, followerId, followeeId int64) error {
	return i.interactiveRepository.CancelFollow(ctx, followerId, followeeId)
}

// IncreaseReadCount 增加阅读数
func (i *InteractiveServiceV1) IncreaseReadCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	return i.interactiveRepository.IncreaseReadCount(ctx, biz, bizId, userId)
}

// IncreaseLikeCount 增加点赞数
func (i *InteractiveServiceV1) IncreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	return i.interactiveRepository.IncreaseLikeCount(ctx, biz, bizId, userId)
}

// DecreaseLikeCount 减少点赞数
func (i *InteractiveServiceV1) DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	return i.interactiveRepository.DecreaseLikeCount(ctx, biz, bizId, userId)
}

// IncreaseCollectCount 增加收藏数
func (i *InteractiveServiceV1) IncreaseCollectCount(ctx context.Context, biz string,
	bizId int64, userId int64, collectId int64) error {
	return i.interactiveRepository.IncreaseCollectCount(ctx, biz, bizId, userId, collectId)
}

// DecreaseCollectCount 减少收藏数
func (i *InteractiveServiceV1) DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId int64, collectId int64) error {
	return i.interactiveRepository.DecreaseCollectCount(ctx, biz, bizId, userId, collectId)
}

// BatchInteractiveReadCount 批量增加阅读数
func (i *InteractiveServiceV1) BatchInteractiveReadCount(ctx context.Context,
	biz string, bizIds, userIds []int64) error {
	return i.interactiveRepository.BatchInteractiveReadCount(ctx, biz, bizIds, userIds)
}

// GetInteractiveInfo 获取互动信息
func (i *InteractiveServiceV1) GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error) {
	return i.interactiveRepository.GetInteractiveInfoByBizId(ctx, biz, bizId)
}
