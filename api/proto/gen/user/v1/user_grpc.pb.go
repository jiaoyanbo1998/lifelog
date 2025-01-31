// protobuf版本

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/proto/user/v1/user.proto

// go包名 user.v1 == userv1

package userv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserService_RegisterByEmailAndPassword_FullMethodName = "/user.v1.UserService/RegisterByEmailAndPassword"
	UserService_LoginByEmailAndPassword_FullMethodName    = "/user.v1.UserService/LoginByEmailAndPassword"
	UserService_GetUserInfoById_FullMethodName            = "/user.v1.UserService/GetUserInfoById"
	UserService_UpdateUserInfoById_FullMethodName         = "/user.v1.UserService/UpdateUserInfoById"
	UserService_DeleteUserInfoByIds_FullMethodName        = "/user.v1.UserService/DeleteUserInfoByIds"
	UserService_Logout_FullMethodName                     = "/user.v1.UserService/Logout"
	UserService_LoginByPhoneCode_FullMethodName           = "/user.v1.UserService/LoginByPhoneCode"
	UserService_UpdateAvatar_FullMethodName               = "/user.v1.UserService/UpdateAvatar"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 定义服务
//go:generate mockgen -source=./user_grpc.pb.go -package=userGRPCMock -destination=mock/usergRPC.mock.go UserServiceClient
type UserServiceClient interface {
	RegisterByEmailAndPassword(ctx context.Context, in *RegisterByEmailAndPasswordRequest, opts ...grpc.CallOption) (*RegisterByEmailAndPasswordResponse, error)
	LoginByEmailAndPassword(ctx context.Context, in *LoginByEmailAndPasswordRequest, opts ...grpc.CallOption) (*LoginByEmailAndPasswordResponse, error)
	GetUserInfoById(ctx context.Context, in *GetUserInfoByIdRequest, opts ...grpc.CallOption) (*GetUserInfoByIdResponse, error)
	UpdateUserInfoById(ctx context.Context, in *UpdateUserInfoByIdRequest, opts ...grpc.CallOption) (*UpdateUserInfoByIdResponse, error)
	DeleteUserInfoByIds(ctx context.Context, in *DeleteUserInfoByIdsRequest, opts ...grpc.CallOption) (*DeleteUserInfoByIdsResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	LoginByPhoneCode(ctx context.Context, in *LoginByPhoneCodeRequest, opts ...grpc.CallOption) (*LoginByPhoneCodeResponse, error)
	UpdateAvatar(ctx context.Context, in *UpdateAvatarRequest, opts ...grpc.CallOption) (*UpdateAvatarResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) RegisterByEmailAndPassword(ctx context.Context, in *RegisterByEmailAndPasswordRequest, opts ...grpc.CallOption) (*RegisterByEmailAndPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterByEmailAndPasswordResponse)
	err := c.cc.Invoke(ctx, UserService_RegisterByEmailAndPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) LoginByEmailAndPassword(ctx context.Context, in *LoginByEmailAndPasswordRequest, opts ...grpc.CallOption) (*LoginByEmailAndPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginByEmailAndPasswordResponse)
	err := c.cc.Invoke(ctx, UserService_LoginByEmailAndPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserInfoById(ctx context.Context, in *GetUserInfoByIdRequest, opts ...grpc.CallOption) (*GetUserInfoByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInfoByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetUserInfoById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUserInfoById(ctx context.Context, in *UpdateUserInfoByIdRequest, opts ...grpc.CallOption) (*UpdateUserInfoByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserInfoByIdResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUserInfoById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUserInfoByIds(ctx context.Context, in *DeleteUserInfoByIdsRequest, opts ...grpc.CallOption) (*DeleteUserInfoByIdsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteUserInfoByIdsResponse)
	err := c.cc.Invoke(ctx, UserService_DeleteUserInfoByIds_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, UserService_Logout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) LoginByPhoneCode(ctx context.Context, in *LoginByPhoneCodeRequest, opts ...grpc.CallOption) (*LoginByPhoneCodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginByPhoneCodeResponse)
	err := c.cc.Invoke(ctx, UserService_LoginByPhoneCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateAvatar(ctx context.Context, in *UpdateAvatarRequest, opts ...grpc.CallOption) (*UpdateAvatarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAvatarResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateAvatar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
//
// 定义服务
type UserServiceServer interface {
	RegisterByEmailAndPassword(context.Context, *RegisterByEmailAndPasswordRequest) (*RegisterByEmailAndPasswordResponse, error)
	LoginByEmailAndPassword(context.Context, *LoginByEmailAndPasswordRequest) (*LoginByEmailAndPasswordResponse, error)
	GetUserInfoById(context.Context, *GetUserInfoByIdRequest) (*GetUserInfoByIdResponse, error)
	UpdateUserInfoById(context.Context, *UpdateUserInfoByIdRequest) (*UpdateUserInfoByIdResponse, error)
	DeleteUserInfoByIds(context.Context, *DeleteUserInfoByIdsRequest) (*DeleteUserInfoByIdsResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	LoginByPhoneCode(context.Context, *LoginByPhoneCodeRequest) (*LoginByPhoneCodeResponse, error)
	UpdateAvatar(context.Context, *UpdateAvatarRequest) (*UpdateAvatarResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) RegisterByEmailAndPassword(context.Context, *RegisterByEmailAndPasswordRequest) (*RegisterByEmailAndPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterByEmailAndPassword not implemented")
}
func (UnimplementedUserServiceServer) LoginByEmailAndPassword(context.Context, *LoginByEmailAndPasswordRequest) (*LoginByEmailAndPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByEmailAndPassword not implemented")
}
func (UnimplementedUserServiceServer) GetUserInfoById(context.Context, *GetUserInfoByIdRequest) (*GetUserInfoByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoById not implemented")
}
func (UnimplementedUserServiceServer) UpdateUserInfoById(context.Context, *UpdateUserInfoByIdRequest) (*UpdateUserInfoByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserInfoById not implemented")
}
func (UnimplementedUserServiceServer) DeleteUserInfoByIds(context.Context, *DeleteUserInfoByIdsRequest) (*DeleteUserInfoByIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserInfoByIds not implemented")
}
func (UnimplementedUserServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedUserServiceServer) LoginByPhoneCode(context.Context, *LoginByPhoneCodeRequest) (*LoginByPhoneCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByPhoneCode not implemented")
}
func (UnimplementedUserServiceServer) UpdateAvatar(context.Context, *UpdateAvatarRequest) (*UpdateAvatarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAvatar not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_RegisterByEmailAndPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterByEmailAndPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RegisterByEmailAndPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RegisterByEmailAndPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RegisterByEmailAndPassword(ctx, req.(*RegisterByEmailAndPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_LoginByEmailAndPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByEmailAndPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginByEmailAndPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_LoginByEmailAndPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginByEmailAndPassword(ctx, req.(*LoginByEmailAndPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserInfoById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserInfoById(ctx, req.(*GetUserInfoByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUserInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserInfoByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUserInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateUserInfoById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUserInfoById(ctx, req.(*UpdateUserInfoByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUserInfoByIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserInfoByIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUserInfoByIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteUserInfoByIds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUserInfoByIds(ctx, req.(*DeleteUserInfoByIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_LoginByPhoneCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByPhoneCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginByPhoneCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_LoginByPhoneCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginByPhoneCode(ctx, req.(*LoginByPhoneCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAvatarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateAvatar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateAvatar(ctx, req.(*UpdateAvatarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterByEmailAndPassword",
			Handler:    _UserService_RegisterByEmailAndPassword_Handler,
		},
		{
			MethodName: "LoginByEmailAndPassword",
			Handler:    _UserService_LoginByEmailAndPassword_Handler,
		},
		{
			MethodName: "GetUserInfoById",
			Handler:    _UserService_GetUserInfoById_Handler,
		},
		{
			MethodName: "UpdateUserInfoById",
			Handler:    _UserService_UpdateUserInfoById_Handler,
		},
		{
			MethodName: "DeleteUserInfoByIds",
			Handler:    _UserService_DeleteUserInfoByIds_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _UserService_Logout_Handler,
		},
		{
			MethodName: "LoginByPhoneCode",
			Handler:    _UserService_LoginByPhoneCode_Handler,
		},
		{
			MethodName: "UpdateAvatar",
			Handler:    _UserService_UpdateAvatar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/user/v1/user.proto",
}
