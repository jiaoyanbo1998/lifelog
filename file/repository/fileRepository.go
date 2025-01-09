package repository

import (
	"context"
	"lifelog-grpc/file/domain"
	"lifelog-grpc/file/repository/dao"
)

type FileRepository interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error)
}

type FileRepositoryV1 struct {
	fileDao dao.FileDao
}

func NewFileService(fileDao dao.FileDao) FileRepository {
	return &FileRepositoryV1{
		fileDao: fileDao,
	}
}
func (f *FileRepositoryV1) GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error) {
	return f.fileDao.GetFile(ctx, fileDomain)
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
