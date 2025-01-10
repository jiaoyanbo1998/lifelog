package repository

import (
	"context"
	"lifelog-grpc/files/domain"
	"lifelog-grpc/files/repository/dao"
)

type FileRepository interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error)
	GetFileByName(ctx context.Context, name string) (domain.FileDomain, error)
}

type FileRepositoryV1 struct {
	fileDao dao.FileDao
}

func NewFileService(fileDao dao.FileDao) FileRepository {
	return &FileRepositoryV1{
		fileDao: fileDao,
	}
}

func (f *FileRepositoryV1) GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error) {
	return f.fileDao.GetFileByUserId(ctx, limit, offset, userId)
}

func (f *FileRepositoryV1) GetFileByName(ctx context.Context, name string) (domain.FileDomain, error) {
	return f.fileDao.GetFileByName(ctx, name)
}

func (f *FileRepositoryV1) CreateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileDao.CreateFile(ctx, fileDomain)
}

func (f *FileRepositoryV1) UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileDao.UpdateFile(ctx, fileDomain)
}

func (f *FileRepositoryV1) DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileDao.DeleteFile(ctx, fileDomain)
}
