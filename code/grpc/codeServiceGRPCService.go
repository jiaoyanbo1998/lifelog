package grpc

import (
	"context"
	"lifelog-grpc/api/proto/gen/code/v1"
	"lifelog-grpc/code/service"
)

// CodeServiceGRPCService 短信服务的grpc服务器
type CodeServiceGRPCService struct {
	codeService service.CodeService
	codev1.UnimplementedCodeServiceServer
}

func NewCodeServiceGRPCService(codeService service.CodeService) *CodeServiceGRPCService {
	return &CodeServiceGRPCService{
		codeService: codeService,
	}
}

func (c *CodeServiceGRPCService) SendPhoneCode(ctx context.Context, request *codev1.SendPhoneCodeRequest) (*codev1.SendPhoneCodeResponse, error) {
	err := c.codeService.SendPhoneCode(ctx, request.GetPhone(), request.GetBiz())
	if err != nil {
		return &codev1.SendPhoneCodeResponse{}, err
	}
	return &codev1.SendPhoneCodeResponse{}, nil
}

func (c *CodeServiceGRPCService) VerifyPhoneCode(ctx context.Context, request *codev1.VerifyPhoneCodeRequest) (*codev1.VerifyPhoneCodeResponse, error) {
	err := c.codeService.VerifyPhoneCode(ctx, request.GetPhone(), request.GetCode(), request.GetBiz())
	if err != nil {
		return &codev1.VerifyPhoneCodeResponse{}, err
	}
	return &codev1.VerifyPhoneCodeResponse{}, nil
}

func (c *CodeServiceGRPCService) SetBlackPhone(ctx context.Context, request *codev1.SetBlackPhoneRequest) (*codev1.SetBlackPhoneResponse, error) {
	err := c.codeService.SetBlackPhone(ctx, request.GetPhone())
	if err != nil {
		return &codev1.SetBlackPhoneResponse{}, err
	}
	return &codev1.SetBlackPhoneResponse{}, nil
}

func (c *CodeServiceGRPCService) IsBackPhone(ctx context.Context, request *codev1.IsBackPhoneRequest) (*codev1.IsBackPhoneResponse, error) {
	isBack, err := c.codeService.IsBackPhone(ctx, request.GetPhone())
	if err != nil {
		return &codev1.IsBackPhoneResponse{}, err
	}
	return &codev1.IsBackPhoneResponse{
		IsBack: isBack,
	}, nil
}
