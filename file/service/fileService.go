package service

import "lifelog-grpc/file/repository"

type FileService interface {
}

type FileServiceV1 struct {
	fileRepository repository.FileRepository
}

func NewFileService(fileRepository repository.FileRepository) FileService {
	return &FileServiceV1{
		fileRepository: fileRepository,
	}
}
