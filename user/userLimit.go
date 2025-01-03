package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
	"lifelog-grpc/pkg/limitx"
)

// UserLimit 装饰器模式，UserService的限流拦截器
type UserLimit struct {
	userv1.UserServiceServer                // UserService的grpc服务
	limiter                  limitx.Limiter // 限流器对象
}

// GetUserInfoById 处理获取用户信息的请求
func (u *UserLimit) GetUserInfoById(ctx context.Context, request *userv1.GetUserInfoByIdRequest) (*userv1.GetUserInfoByIdResponse, error) {
	// 生成限流key
	key := fmt.Sprintf("limiter:user:GetUserInfoById:%d", request.GetUserDomain().GetId())
	// 调用限流器，判断是否需要限流
	limit, err := u.limiter.Limit(ctx, key)
	// 触发限流失败
	if err != nil {
		return nil, errors.New("触发限流失败")
	}
	// 触发限流
	if limit == true {
		return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
	}
	// 执行下一个拦截器，或者是真实的业务代码
	return u.UserServiceServer.GetUserInfoById(ctx, request)
}
