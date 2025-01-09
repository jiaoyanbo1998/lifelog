package grpc

import "lifelog-grpc/file/service"

type FileServiceGRPCService struct {
	fileService service.FileService
}

func NewFileServiceGRPCService(fileService service.FileService) *FileServiceGRPCService {
	return &FileServiceGRPCService{
		fileService: fileService,
	}
}
