package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"lifelog-grpc/files/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type FileDao interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error)
	GetFileByName(ctx context.Context, name string) (domain.FileDomain, error)
}

type FileDaoV1 struct {
	db     *gorm.DB
	logger loggerx.Logger
}

func NewFileDao(db *gorm.DB, logger loggerx.Logger) FileDao {
	return &FileDaoV1{
		db:     db,
		logger: logger,
	}
}

type File struct {
	Id         int64
	UserId     int64
	Url        string
	Name       string
	Content    string
	CreateTime int64
	UpdateTime int64
}

func (File) TableName() string {
	return "tb_file"
}

func (f *FileDaoV1) GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error) {
	var files []File
	err := f.db.WithContext(ctx).Where("user_id = ?", userId).
		Limit(int(limit)).Offset(int(offset)).Find(&files).Error
	if err != nil {
		f.logger.Error("根据用户id获取文件信息失败", loggerx.Error(err),
			loggerx.Int("userId", int(userId)))
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("未找到文件")
	}
	// []File，转换为[]domain.FileDomain
	return convertToFileDomains(files), nil
}

func (f *FileDaoV1) GetFileByName(ctx context.Context, name string) (domain.FileDomain, error) {
	var file File
	err := f.db.WithContext(ctx).Where("name = ?", name).First(&file).Error
	if err != nil {
		// 记录未找到
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.FileDomain{}, errors.New("未找到文件")
		}
		// 其他错误
		return domain.FileDomain{}, fmt.Errorf("未知错误：%w", err)
	}
	// 将File，转换为domain.FileDomain
	return convertToFileDomain(file), nil
}

func (f *FileDaoV1) CreateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	now := time.Now().UnixMilli()
	file := File{
		UserId:     fileDomain.UserId,
		Url:        fileDomain.Url,
		Name:       fileDomain.Name,
		Content:    fileDomain.Content,
		CreateTime: now,
		UpdateTime: now,
	}
	err := f.db.WithContext(ctx).Create(&file).Error
	if err != nil {
		f.logger.Error("将文件信息插入数据库失败", loggerx.Error(err),
			loggerx.Int("userId", int(fileDomain.UserId)),
			loggerx.String("name", fileDomain.Name),
			loggerx.String("url", fileDomain.Url))
		return err
	}
	return nil
}

func (f *FileDaoV1) UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	now := time.Now().UnixMilli()
	tx := f.db.WithContext(ctx).Model(&File{}).
		Where("user_id = ? and id = ?", fileDomain.UserId, fileDomain.Id).
		Updates(map[string]interface{}{
			"url":         fileDomain.Url,
			"name":        fileDomain.Name,
			"content":     fileDomain.Content,
			"update_time": now,
		})
	if tx.Error != nil {
		f.logger.Error("将文件信息更新数据库失败", loggerx.Error(tx.Error),
			loggerx.Int("userId", int(fileDomain.UserId)),
			loggerx.String("name", fileDomain.Name),
			loggerx.String("url", fileDomain.Url))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("未找到文件")
	}
	return nil
}

func (f *FileDaoV1) DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error {
	err := f.db.WithContext(ctx).
		Where("user_id = ? and id = ?", fileDomain.UserId, fileDomain.Id).
		Delete(&File{}).Error
	if err != nil {
		f.logger.Error("将文件信息删除数据库失败", loggerx.Error(err),
			loggerx.Int("userId", int(fileDomain.UserId)),
			loggerx.Int("id", int(fileDomain.Id)))
		return err
	}
	return nil
}

// 将File，转换为domain.FileDomain
func convertToFileDomain(file File) domain.FileDomain {
	return domain.FileDomain{
		Id:         file.Id,
		UserId:     file.UserId,
		Url:        file.Url,
		Name:       file.Name,
		Content:    file.Content,
		CreateTime: file.CreateTime,
		UpdateTime: file.UpdateTime,
	}
}

// 将[]File，转换为[]domain.FileDomain
func convertToFileDomains(files []File) []domain.FileDomain {
	fds := make([]domain.FileDomain, 0, len(files))
	for _, file := range files {
		fds = append(fds, convertToFileDomain(file))
	}
	return fds
}
