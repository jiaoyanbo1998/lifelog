package grpc

import (
	"context"
	codev1 "lifelog-grpc/api/proto/gen/code/v1"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/user/domain"
	"lifelog-grpc/user/service"
)

// UserServiceGRPCService User服务的grpc服务器，在这一层聚合不同的服务
type UserServiceGRPCService struct {
	userService       service.UserService
	logger            loggerx.Logger
	codeServiceClient codev1.CodeServiceClient
	userv1.UnimplementedUserServiceServer
}

func NewUserServiceGRPCService(userService service.UserService,
	l loggerx.Logger,
	codeServiceClient codev1.CodeServiceClient) *UserServiceGRPCService {
	return &UserServiceGRPCService{
		userService:       userService,
		logger:            l,
		codeServiceClient: codeServiceClient,
	}
}

func (u *UserServiceGRPCService) UpdateAvatar(ctx context.Context, request *userv1.UpdateAvatarRequest) (*userv1.UpdateAvatarResponse, error) {
	err := u.userService.UpdateAvatar(ctx, request.GetUserId(), request.GetFilePath())
	if err != nil {
		return &userv1.UpdateAvatarResponse{}, err
	}
	return &userv1.UpdateAvatarResponse{}, nil
}

// RegisterByEmailAndPassword 邮箱密码注册
func (u *UserServiceGRPCService) RegisterByEmailAndPassword(ctx context.Context, request *userv1.RegisterByEmailAndPasswordRequest) (*userv1.RegisterByEmailAndPasswordResponse, error) {
	// 调用service层
	res, err := u.userService.RegisterByEmailAndPassword(ctx, domain.UserDomain{
		Email:    request.GetUserDomain().Email,
		Password: request.GetUserDomain().Password,
	})
	if err != nil {
		return &userv1.RegisterByEmailAndPasswordResponse{}, err
	}
	return &userv1.RegisterByEmailAndPasswordResponse{
		UserDomain: &userv1.UserDomain{
			Id:       res.Id,
			Email:    res.Email,
			NickName: res.NickName,
		},
	}, nil
}

// LoginByEmailAndPassword 邮箱密码登录
func (u *UserServiceGRPCService) LoginByEmailAndPassword(ctx context.Context, request *userv1.LoginByEmailAndPasswordRequest) (*userv1.LoginByEmailAndPasswordResponse, error) {
	// 调用service层
	res, info, err := u.userService.LoginByEmailAndPassword(ctx, domain.UserDomain{
		Email:    request.GetUserDomain().Email,
		Password: request.GetUserDomain().Password,
	})
	if err != nil {
		return &userv1.LoginByEmailAndPasswordResponse{}, err
	}
	return &userv1.LoginByEmailAndPasswordResponse{
		UserDomain: &userv1.UserDomain{
			Id:        res.Id,
			NickName:  res.NickName,
			Authority: res.Authority,
		},
		Info: info,
	}, nil
}

// GetUserInfoById 根据id获取用户信息
func (u *UserServiceGRPCService) GetUserInfoById(ctx context.Context, request *userv1.GetUserInfoByIdRequest) (*userv1.GetUserInfoByIdResponse, error) {
	// 调用service层
	res, err := u.userService.GetUserInfoById(ctx, domain.UserDomain{
		Id: request.GetUserDomain().Id,
	})
	if err != nil {
		return &userv1.GetUserInfoByIdResponse{}, err
	}
	return &userv1.GetUserInfoByIdResponse{
		UserDomain: &userv1.UserDomain{
			Id:       res.Id,
			Email:    res.Email,
			Phone:    res.Phone,
			NickName: res.NickName,
		},
	}, nil
}

// UpdateUserInfoById 根据用户id更新用户信息
func (u *UserServiceGRPCService) UpdateUserInfoById(ctx context.Context, request *userv1.UpdateUserInfoByIdRequest) (*userv1.UpdateUserInfoByIdResponse, error) {
	// 调用service层
	err := u.userService.UpdateUserInfoById(ctx, domain.UserDomain{
		Id:          request.GetUserDomain().Id,
		Password:    request.GetUserDomain().Password,
		NickName:    request.GetUserDomain().NickName,
		Phone:       request.GetUserDomain().Phone,
		Email:       request.GetUserDomain().Email,
		NewPassword: request.GetUserDomain().NewPassword,
	})
	if err != nil {
		return &userv1.UpdateUserInfoByIdResponse{}, err
	}
	return &userv1.UpdateUserInfoByIdResponse{}, nil
}

// DeleteUserInfoByIds 根据用户id删除用户信息
func (u *UserServiceGRPCService) DeleteUserInfoByIds(ctx context.Context, request *userv1.DeleteUserInfoByIdsRequest) (*userv1.DeleteUserInfoByIdsResponse, error) {
	// 调用service层
	err := u.userService.DeleteUSerInfoByIds(ctx, request.GetIds())
	if err != nil {
		return &userv1.DeleteUserInfoByIdsResponse{}, err
	}
	return &userv1.DeleteUserInfoByIdsResponse{}, nil
}

// Logout 退出
func (u *UserServiceGRPCService) Logout(ctx context.Context, request *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	// 调用service层
	err := u.userService.Logout(ctx, request.GetSessionId())
	if err != nil {
		return &userv1.LogoutResponse{}, err
	}
	return &userv1.LogoutResponse{}, nil
}

// LoginByPhoneCode 手机号验证码登录
func (u *UserServiceGRPCService) LoginByPhoneCode(ctx context.Context, request *userv1.LoginByPhoneCodeRequest) (*userv1.LoginByPhoneCodeResponse, error) {
	// 调用service层
	res, err := u.userService.LoginByPhoneCode(ctx, domain.UserDomain{
		Phone: request.GetUserDomain().Phone,
		Code:  request.GetUserDomain().Code,
	}, request.Biz)
	if err != nil {
		return &userv1.LoginByPhoneCodeResponse{}, err
	}
	// 验证短信验证码是否正确
	_, err = u.codeServiceClient.VerifyPhoneCode(ctx, &codev1.VerifyPhoneCodeRequest{
		Phone: request.GetUserDomain().Phone,
		Code:  request.GetUserDomain().Code,
		Biz:   request.GetBiz(),
	})
	// 验证失败
	if err != nil {
		return &userv1.LoginByPhoneCodeResponse{}, err
	}
	// 验证成功
	return &userv1.LoginByPhoneCodeResponse{
		UserDomain: &userv1.UserDomain{
			Id:        res.Id,
			NickName:  res.NickName,
			Authority: res.Authority,
		},
	}, nil
}
