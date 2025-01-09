package grpc

import (
	"context"
	filev1 "lifelog-grpc/api/proto/gen/file/v1"
	"lifelog-grpc/file/domain"
	"lifelog-grpc/file/service"
)

type FileServiceGRPCService struct {
	fileService service.FileService
	filev1.UnimplementedFileServiceServer
}

func NewFileServiceGRPCService(fileService service.FileService) *FileServiceGRPCService {
	return &FileServiceGRPCService{
		fileService: fileService,
	}
}

func (f *FileServiceGRPCService) GetFile(ctx context.Context, request *filev1.GetFileRequest) (*filev1.GetFileResponse, error) {
	err := f.fileService.CreateFile(ctx, domain.FileDomain{
		Biz:   request.GetFile().GetBiz(),
		BizId: request.GetFile().GetBizId(),
	})
	if err != nil {
		return &filev1.GetFileResponse{}, err
	}
	return &filev1.GetFileResponse{}, nil
}

func (f *FileServiceGRPCService) InsertFile(ctx context.Context, request *filev1.InsertFileRequest) (*filev1.InsertFileResponse, error) {
	err := f.fileService.CreateFile(ctx, domain.FileDomain{
		Biz:   request.GetFile().GetBiz(),
		BizId: request.GetFile().GetBizId(),
		Url:   request.GetFile().GetUrl(),
	})
	if err != nil {
		return &filev1.InsertFileResponse{}, err
	}
	return &filev1.InsertFileResponse{}, nil
}

func (f *FileServiceGRPCService) UpdateFile(ctx context.Context, request *filev1.UpdateFileRequest) (*filev1.UpdateFileResponse, error) {
	f.fileService.UpdateFile(ctx, domain.FileDomain{
		Biz:   request.GetFile().GetBiz(),
		BizId: request.GetFile().GetBizId(),
		Url:   request.GetFile().GetUrl(),
	})
	return &filev1.UpdateFileResponse{}, nil
}

func (f *FileServiceGRPCService) DeleteFile(ctx context.Context, request *filev1.DeleteFileRequest) (*filev1.DeleteFileResponse, error) {
	err := f.fileService.DeleteFile(ctx, domain.FileDomain{
		Biz:   request.GetFile().GetBiz(),
		BizId: request.GetFile().GetBizId(),
	})
	if err != nil {
		return &filev1.DeleteFileResponse{}, err
	}
	return &filev1.DeleteFileResponse{}, nil
}
