package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"lifelog-grpc/file/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type FileDao interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error)
}

type FileDaoV1 struct {
	db     *gorm.DB
	logger loggerx.Logger
}

func NewFileService(db *gorm.DB, logger loggerx.Logger) FileDao {
	return &FileDaoV1{
		db:     db,
		logger: logger,
	}
}

type File struct {
	Id         int64
	Biz        string
	BizId      int64
	Url        string
	CreateTime int64
	UpdateTime int64
}

func (File) TableName() string {
	return "tb_file"
}

func (f *FileDaoV1) GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error) {
	var files []File
	res := f.db.WithContext(ctx).Where("biz = ? and biz_id = ?",
		fileDomain.Biz, fileDomain.BizId).Find(&files)
	if res.RowsAffected == 0 {
		return nil, errors.New("未找到文件")
	}
	var urls []string
	for _, file := range files {
		urls = append(urls, file.Url)
	}
	return urls, nil
}

func (f *FileDaoV1) CreateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	now := time.Now().UnixMilli()
	var file File
	file.Biz = fileDomain.Biz
	file.BizId = fileDomain.BizId
	file.Url = fileDomain.Url
	file.CreateTime = now
	file.UpdateTime = now
	err := f.db.WithContext(ctx).Create(&file).Error
	if err != nil {
		f.logger.Error("将文件信息，插入数据库失败", loggerx.Error(err),
			loggerx.String("url", fileDomain.Url),
			loggerx.Int("userId", int(fileDomain.BizId)),
			loggerx.String("biz", fileDomain.Biz))
		return err
	}
	return nil
}

func (f *FileDaoV1) UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	now := time.Now().UnixMilli()
	err := f.db.WithContext(ctx).Model(&File{}).
		Where("biz_id = ? and biz = ?", fileDomain.BizId, fileDomain.Biz).
		Updates(map[string]interface{}{
			"url":         fileDomain.Url,
			"update_time": now,
		}).Error
	if err != nil {
		f.logger.Error("将文件信息，更新数据库失败", loggerx.Error(err),
			loggerx.String("url", fileDomain.Url),
			loggerx.Int("userId", int(fileDomain.BizId)),
			loggerx.String("biz", fileDomain.Biz))
		return err
	}
	return nil
}

func (f *FileDaoV1) DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error {
	err := f.db.WithContext(ctx).
		Where("biz_id =? and biz =?", fileDomain.BizId, fileDomain.Biz).
		Delete(&File{}).Error
	if err != nil {
		f.logger.Error("将文件信息，删除数据库失败", loggerx.Error(err),
			loggerx.String("url", fileDomain.Url),
			loggerx.Int("userId", int(fileDomain.BizId)),
			loggerx.String("biz", fileDomain.Biz))
		return err
	}
	return nil
}
