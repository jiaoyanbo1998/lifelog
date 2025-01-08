package repository

import (
	"context"
	"database/sql"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/repository/dao"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, commentDomain domain.CommentDomain) error
	DeleteComment(ctx context.Context, id int64) error
	EditComment(ctx context.Context, commentDomain domain.CommentDomain) error
	FirstList(ctx context.Context, biz string, lifeLogId int64, min int64) ([]domain.CommentDomain, error)
	EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error)
	SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error)
	BatchCreateComment(ctx context.Context, comments []domain.CommentDomain) error
}

type CommentRepositoryV1 struct {
	commentDao dao.CommentDao
}

func NewCommentRepository(commentDao dao.CommentDao) CommentRepository {
	return &CommentRepositoryV1{
		commentDao: commentDao,
	}
}

func (c *CommentRepositoryV1) CreateComment(ctx context.Context,
	commentDomain domain.CommentDomain) error {
	return c.commentDao.InsertComment(ctx, dao.Comment{
		Biz:     commentDomain.Biz,
		BizId:   commentDomain.BizId,
		Content: commentDomain.Content,
		UserId:  commentDomain.UserId,
		ParentId: sql.NullInt64{
			Int64: commentDomain.ParentId,
			Valid: commentDomain.ParentId != 0,
		},
		RootId: sql.NullInt64{
			Int64: commentDomain.RootId,
			Valid: commentDomain.RootId != 0,
		},
		Uuid: commentDomain.Uuid,
	})
}

func (c *CommentRepositoryV1) BatchCreateComment(ctx context.Context, comments []domain.CommentDomain) error {
	return c.commentDao.BatchInsertComment(ctx, comments)
}

func (c *CommentRepositoryV1) DeleteComment(ctx context.Context, id int64) error {
	return c.commentDao.DeleteComment(ctx, id)
}

func (c *CommentRepositoryV1) EditComment(ctx context.Context, commentDomain domain.CommentDomain) error {
	return c.commentDao.UpdateComment(ctx, dao.Comment{
		Id:      commentDomain.Id,
		Content: commentDomain.Content,
	})
}

func (c *CommentRepositoryV1) FirstList(ctx context.Context, biz string, lifeLogId int64,
	min int64) ([]domain.CommentDomain, error) {
	return c.commentDao.FirstList(ctx, biz, lifeLogId, min)
}

func (c *CommentRepositoryV1) EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error) {
	return c.commentDao.EveryRootChildSonList(ctx, id, RootId, limit)
}

func (c *CommentRepositoryV1) SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error) {
	return c.commentDao.SonList(ctx, parentId, limit, offset)
}
