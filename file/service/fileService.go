package service

import (
	"context"
	"lifelog-grpc/file/domain"
	"lifelog-grpc/file/repository"
	"lifelog-grpc/pkg/loggerx"
)

type FileService interface {
	CreateFile(ctx context.Context, fileDomain domain.FileDomain) error
	UpdateFile(ctx context.Context, fileDomain domain.FileDomain) error
	DeleteFile(ctx context.Context, fileDomain domain.FileDomain) error
	GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error)
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

func (f *FileServiceV1) GetFile(ctx context.Context, fileDomain domain.FileDomain) ([]string, error) {
	return f.fileRepository.GetFile(ctx, fileDomain)
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
