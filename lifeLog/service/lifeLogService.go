package service

import (
	"context"
	"lifelog-grpc/lifeLog/domain"
	"lifelog-grpc/lifeLog/repository"
)

type LifeLogService interface {
	Save(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error)
	Delete(ctx context.Context, ids []int64, public bool) error
	SearchByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error)
	Detail(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error)
	SearchByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error)
	Revoke(ctx context.Context, id, authorId int64) error
	Publish(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error
	DetailMany(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error)
}

type LifeLogServiceV1 struct {
	lifeLogRepository repository.LifeLogRepository
}

func NewLifeLogService(lifeLogRepository repository.LifeLogRepository) LifeLogService {
	return &LifeLogServiceV1{
		lifeLogRepository: lifeLogRepository,
	}
}

// Save 创建或修改LifeLog
func (a *LifeLogServiceV1) Save(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error) {
	// 传入id表示是修改
	if lifeLogDomain.Id > 0 {
		return a.lifeLogRepository.Modify(ctx, lifeLogDomain)
	}
	// 不传入id表示是创建
	return a.lifeLogRepository.Create(ctx, lifeLogDomain)
}

// Delete 删除LifeLog
func (a *LifeLogServiceV1) Delete(ctx context.Context, ids []int64, public bool) error {
	return a.lifeLogRepository.Delete(ctx, ids, public)
}

// SearchByTitle 查询LifeLog
func (a *LifeLogServiceV1) SearchByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error) {
	return a.lifeLogRepository.SearchByTitle(ctx, title, limit, offset)
}

// Detail 查询LifeLog
func (a *LifeLogServiceV1) Detail(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error) {
	return a.lifeLogRepository.SearchById(ctx, id, public)
}

// DetailMany 查询多个LifeLog
func (a *LifeLogServiceV1) DetailMany(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error) {
	return a.lifeLogRepository.SearchByIds(ctx, ids)
}

// SearchByAuthorId 查询LifeLog
func (a *LifeLogServiceV1) SearchByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error) {
	return a.lifeLogRepository.SearchByAuthorId(ctx, authorId, limit, offset)
}

// Revoke 撤销LifeLog
func (a *LifeLogServiceV1) Revoke(ctx context.Context, id, authorId int64) error {
	return a.lifeLogRepository.RevokeById(ctx, id, authorId)
}

// Publish 发布LifeLog
func (a *LifeLogServiceV1) Publish(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error {
	lifeLogDomain.Status = domain.LifeLogStatusPublished
	return a.lifeLogRepository.Sync(ctx, lifeLogDomain)
}
