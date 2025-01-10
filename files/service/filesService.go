package service

import (
	"context"
	"lifelog-grpc/files/domain"
	"lifelog-grpc/files/repository"
	"lifelog-grpc/pkg/loggerx"
)

type FileService interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error)
	GetFileByName(ctx context.Context, name string) (domain.FileDomain, error)
}

type FileServiceV1 struct {
	fileRepository repository.FileRepository
	logger         loggerx.Logger
}

func NewFileService(fileRepository repository.FileRepository, logger loggerx.Logger) FileService {
	return &FileServiceV1{
		fileRepository: fileRepository,
		logger:         logger,
	}
}

func (f *FileServiceV1) GetFileByUserId(ctx context.Context, limit, offset, userId int64) ([]domain.FileDomain, error) {
	return f.fileRepository.GetFileByUserId(ctx, limit, offset, userId)
}

func (f *FileServiceV1) GetFileByName(ctx context.Context, name string) (domain.FileDomain, error) {
	return f.fileRepository.GetFileByName(ctx, name)
}

func (f *FileServiceV1) CreateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileRepository.CreateFile(ctx, fileDomain)
}

func (f *FileServiceV1) UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileRepository.UpdateFile(ctx, fileDomain)
}

func (f *FileServiceV1) DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error {
	return f.fileRepository.DeleteFile(ctx, fileDomain)
}
