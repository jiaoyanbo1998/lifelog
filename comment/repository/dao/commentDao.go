package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CommentDao interface {
	InsertComment(ctx context.Context, comment Comment) error
	DeleteComment(ctx context.Context, id int64) error
	UpdateComment(ctx context.Context, comment Comment) error
	FirstList(ctx context.Context, biz string, lifeLogId int64, min int64) ([]domain.CommentDomain, error)
	EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error)
	SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error)
	BatchInsertComment(ctx context.Context, comments []domain.CommentDomain) error
}

type CommentDaoGorm struct {
	db     *gorm.DB
	logger loggerx.Logger
}

func NewCommentDaoGorm(db *gorm.DB, logger loggerx.Logger) CommentDao {
	return &CommentDaoGorm{
		db:     db,
		logger: logger,
	}
}

type Comment struct {
	Id           int64         `json:"id"`
	UserId       int64         `json:"user_id"`
	Biz          string        `json:"biz"`
	BizId        int64         `json:"biz_id"`
	Content      string        `json:"content"`
	ParentId     sql.NullInt64 `json:"parent_id"`
	RootId       sql.NullInt64 `json:"root_id"`
	CreateTime   int64         `json:"create_time"`
	UpdateTime   int64         `json:"update_time"`
	Uuid         string        `json:"uuid"`
	TargetUserId int64         `json:"target_user_id"`
}

func (Comment) TableName() string {
	return "tb_comment"
}

func (c *CommentDaoGorm) BatchInsertComment(ctx context.Context, commentdomains []domain.CommentDomain) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comments []Comment
		now := time.Now().UnixMilli()
		for _, comment := range commentdomains {
			comments = append(comments, Comment{
				UserId:  comment.UserId,
				Biz:     comment.Biz,
				BizId:   comment.BizId,
				Content: comment.Content,
				ParentId: sql.NullInt64{
					Int64: comment.ParentId,
					Valid: comment.ParentId != 0,
				},
				RootId: sql.NullInt64{
					Int64: comment.RootId,
					Valid: comment.RootId != 0,
				},
				CreateTime:   now,
				UpdateTime:   now,
				Uuid:         comment.Uuid,
				TargetUserId: comment.TargetUserId,
			})
		}
		return c.db.WithContext(ctx).Create(&comments).Error
	})
}

func (c *CommentDaoGorm) InsertComment(ctx context.Context, comment Comment) error {
	now := time.Now().UnixMilli()
	comment.CreateTime = now
	comment.UpdateTime = now
	return c.db.WithContext(ctx).Create(&comment).Error
}

func (c *CommentDaoGorm) DeleteComment(ctx context.Context, id int64) error {
	return c.db.WithContext(ctx).Delete(&Comment{
		Id: id,
	}).Error
}

func (c *CommentDaoGorm) UpdateComment(ctx context.Context, comment Comment) error {
	now := time.Now().UnixMilli()
	return c.db.WithContext(ctx).Where("id = ?", comment.Id).
		Model(&Comment{}).Updates(map[string]any{
		"content":     comment.Content,
		"update_time": now,
	}).Error
}

func (c *CommentDaoGorm) FirstList(ctx context.Context, biz string, lifeLogId int64,
	min int64) ([]domain.CommentDomain, error) {
	var comments []Comment
	// parent_id == null 代表根结点
	//  查询出所有根结点
	err := c.db.WithContext(ctx).Where(
		"biz = ? and biz_id = ? and parent_id is null",
		biz, lifeLogId).
		Limit(int(min)).
		Find(&comments).
		Error
	if len(comments) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return c.commentToCommentDomains(comments), nil
}

func (c *CommentDaoGorm) commentToCommentDomains(comments []Comment) []domain.CommentDomain {
	var domains []domain.CommentDomain
	for _, comment := range comments {
		domains = append(domains, domain.CommentDomain{
			Id:      comment.Id,
			UserId:  comment.UserId,
			Content: comment.Content,
		})
	}
	return domains
}

func (c *CommentDaoGorm) EveryRootChildSonList(ctx context.Context, id, RootId, limit int64) ([]domain.CommentDomain, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", RootId, id). // 游标，此页的最后一个数据的id，用于提高分页查询的效率
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	if len(res) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return c.commentToCommentDomains(res), nil
}

func (c *CommentDaoGorm) SonList(ctx context.Context, parentId int64, limit int64, offset int64) ([]domain.CommentDomain, error) {
	var comments []Comment
	err := c.db.WithContext(ctx).Where("parent_id = ?", parentId).Find(&comments).
		Offset(int(offset)).Limit(int(limit)).Error
	if len(comments) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return c.commentToCommentDomains(comments), nil
}
