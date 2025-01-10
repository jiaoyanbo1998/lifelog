package grpc

import (
	"context"
	filesv1 "lifelog-grpc/api/proto/gen/files/v1"
	"lifelog-grpc/files/domain"
	"lifelog-grpc/files/service"
)

type FilesServiceGRPCService struct {
	fileService service.FileService
	filesv1.UnimplementedFilesServiceServer
}

func NewFileServiceGRPCService(fileService service.FileService) *FilesServiceGRPCService {
	return &FilesServiceGRPCService{
		fileService: fileService,
	}
}

func (f *FilesServiceGRPCService) CreateFile(ctx context.Context, request *filesv1.CreateFileRequest) (*filesv1.CreateFileResponse, error) {
	err := f.fileService.CreateFile(ctx, domain.FileDomain{
		UserId:  request.GetFile().GetUserId(),
		Url:     request.GetFile().GetUrl(),
		Name:    request.GetFile().GetName(),
		Content: request.GetFile().GetContent(),
	})
	if err != nil {
		return &filesv1.CreateFileResponse{}, err
	}
	return &filesv1.CreateFileResponse{}, nil
}

func (f *FilesServiceGRPCService) UpdateFile(ctx context.Context, request *filesv1.UpdateFileRequest) (*filesv1.UpdateFileResponse, error) {
	err := f.fileService.UpdateFile(ctx, domain.FileDomain{
		Id:      request.GetFile().GetId(),
		UserId:  request.GetFile().GetUserId(),
		Url:     request.GetFile().GetUrl(),
		Name:    request.GetFile().GetName(),
		Content: request.GetFile().GetContent(),
	})
	if err != nil {
		return &filesv1.UpdateFileResponse{}, err
	}
	return &filesv1.UpdateFileResponse{}, nil
}

func (f *FilesServiceGRPCService) DeleteFile(ctx context.Context, request *filesv1.DeleteFileRequest) (*filesv1.DeleteFileResponse, error) {
	err := f.fileService.DeleteFile(ctx, domain.FileDomain{
		Id:     request.GetFile().GetId(),
		UserId: request.GetFile().GetUserId(),
	})
	if err != nil {
		return &filesv1.DeleteFileResponse{}, err
	}
	return &filesv1.DeleteFileResponse{}, nil
}

func (f *FilesServiceGRPCService) GetFileByUserId(ctx context.Context, request *filesv1.GetFileByUserIdRequest) (*filesv1.GetFileByUserIdResponse, error) {
	files, err := f.fileService.GetFileByUserId(ctx, request.GetLimit(), request.GetOffset(), request.GetFile().UserId)
	if err != nil {
		return &filesv1.GetFileByUserIdResponse{}, err
	}
	fs := make([]*filesv1.File, 0, len(files))
	for _, file := range files {
		fs = append(fs, &filesv1.File{
			Id:         file.Id,
			UserId:     file.UserId,
			Url:        file.Url,
			Name:       file.Name,
			Content:    file.Content,
			CreateTime: file.CreateTime,
			UpdateTime: file.UpdateTime,
		})
	}
	return &filesv1.GetFileByUserIdResponse{
		File: fs,
	}, nil
}

func (f *FilesServiceGRPCService) GetFileByName(ctx context.Context, request *filesv1.GetFileByNameRequest) (*filesv1.GetFileByNameResponse, error) {
	info, err := f.fileService.GetFileByName(ctx, request.GetFile().Name)
	if err != nil {
		return &filesv1.GetFileByNameResponse{}, err
	}
	return &filesv1.GetFileByNameResponse{
		File: &filesv1.File{
			Id:         info.Id,
			UserId:     info.UserId,
			Url:        info.Url,
			Name:       info.Name,
			Content:    info.Content,
			CreateTime: info.CreateTime,
			UpdateTime: info.UpdateTime,
		},
	}, nil
}
