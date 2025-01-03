package service

import (
	"golang.org/x/net/context"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, commentDomain domain.CommentDomain) error
	DeleteComment(ctx context.Context, id int64) error
	EditComment(ctx context.Context, commentDomain domain.CommentDomain) error
	FirstList(ctx context.Context, biz string, lifeLogId int64, min int64) (
		[]domain.CommentDomain, error)
	EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error)
	SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error)
	BatchCreateComment(ctx context.Context, comments []domain.CommentDomain) error
}

type CommentServiceV1 struct {
	repository repository.CommentRepository
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &CommentServiceV1{
		repository: repository,
	}
}

func (c *CommentServiceV1) BatchCreateComment(ctx context.Context, comments []domain.CommentDomain) error {
	return c.repository.BatchCreateComment(ctx, comments)
}

func (c *CommentServiceV1) CreateComment(ctx context.Context,
	commentDomain domain.CommentDomain) error {
	return c.repository.CreateComment(ctx, commentDomain)
}

func (c *CommentServiceV1) DeleteComment(ctx context.Context, id int64) error {
	return c.repository.DeleteComment(ctx, id)
}

func (c *CommentServiceV1) EditComment(ctx context.Context, commentDomain domain.CommentDomain) error {
	return c.repository.EditComment(ctx, commentDomain)
}

func (c *CommentServiceV1) FirstList(ctx context.Context, biz string, lifeLogId int64,
	min int64) ([]domain.CommentDomain, error) {
	return c.repository.FirstList(ctx, biz, lifeLogId, min)
}

func (c *CommentServiceV1) EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error) {
	return c.repository.EveryRootChildSonList(ctx, id, RootId, limit)
}

func (c *CommentServiceV1) SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error) {
	return c.repository.SonList(ctx, parentId, limit, offset)
}
