// 版本

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/proto/lifelog/v1/lifelog.proto

// 包名

package lifelogv1

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
	LifeLogService_Edit_FullMethodName          = "/lifelog.v1.LifeLogService/Edit"
	LifeLogService_Delete_FullMethodName        = "/lifelog.v1.LifeLogService/Delete"
	LifeLogService_SearchByTitle_FullMethodName = "/lifelog.v1.LifeLogService/SearchByTitle"
	LifeLogService_DraftList_FullMethodName     = "/lifelog.v1.LifeLogService/DraftList"
	LifeLogService_Revoke_FullMethodName        = "/lifelog.v1.LifeLogService/Revoke"
	LifeLogService_Publish_FullMethodName       = "/lifelog.v1.LifeLogService/Publish"
	LifeLogService_Detail_FullMethodName        = "/lifelog.v1.LifeLogService/Detail"
)

// LifeLogServiceClient is the client API for LifeLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// lifelog服务
type LifeLogServiceClient interface {
	Edit(ctx context.Context, in *EditRequest, opts ...grpc.CallOption) (*EditResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	SearchByTitle(ctx context.Context, in *SearchByTitleRequest, opts ...grpc.CallOption) (*SearchByTitleResponse, error)
	DraftList(ctx context.Context, in *DraftListRequest, opts ...grpc.CallOption) (*DraftListResponse, error)
	Revoke(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*RevokeResponse, error)
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
	Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error)
}

type lifeLogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLifeLogServiceClient(cc grpc.ClientConnInterface) LifeLogServiceClient {
	return &lifeLogServiceClient{cc}
}

func (c *lifeLogServiceClient) Edit(ctx context.Context, in *EditRequest, opts ...grpc.CallOption) (*EditResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EditResponse)
	err := c.cc.Invoke(ctx, LifeLogService_Edit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, LifeLogService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) SearchByTitle(ctx context.Context, in *SearchByTitleRequest, opts ...grpc.CallOption) (*SearchByTitleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchByTitleResponse)
	err := c.cc.Invoke(ctx, LifeLogService_SearchByTitle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) DraftList(ctx context.Context, in *DraftListRequest, opts ...grpc.CallOption) (*DraftListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DraftListResponse)
	err := c.cc.Invoke(ctx, LifeLogService_DraftList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) Revoke(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*RevokeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeResponse)
	err := c.cc.Invoke(ctx, LifeLogService_Revoke_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, LifeLogService_Publish_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lifeLogServiceClient) Detail(ctx context.Context, in *DetailRequest, opts ...grpc.CallOption) (*DetailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetailResponse)
	err := c.cc.Invoke(ctx, LifeLogService_Detail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LifeLogServiceServer is the server API for LifeLogService service.
// All implementations must embed UnimplementedLifeLogServiceServer
// for forward compatibility.
//
// lifelog服务
type LifeLogServiceServer interface {
	Edit(context.Context, *EditRequest) (*EditResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	SearchByTitle(context.Context, *SearchByTitleRequest) (*SearchByTitleResponse, error)
	DraftList(context.Context, *DraftListRequest) (*DraftListResponse, error)
	Revoke(context.Context, *RevokeRequest) (*RevokeResponse, error)
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	Detail(context.Context, *DetailRequest) (*DetailResponse, error)
	mustEmbedUnimplementedLifeLogServiceServer()
}

// UnimplementedLifeLogServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLifeLogServiceServer struct{}

func (UnimplementedLifeLogServiceServer) Edit(context.Context, *EditRequest) (*EditResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedLifeLogServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedLifeLogServiceServer) SearchByTitle(context.Context, *SearchByTitleRequest) (*SearchByTitleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByTitle not implemented")
}
func (UnimplementedLifeLogServiceServer) DraftList(context.Context, *DraftListRequest) (*DraftListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DraftList not implemented")
}
func (UnimplementedLifeLogServiceServer) Revoke(context.Context, *RevokeRequest) (*RevokeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}
func (UnimplementedLifeLogServiceServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedLifeLogServiceServer) Detail(context.Context, *DetailRequest) (*DetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detail not implemented")
}
func (UnimplementedLifeLogServiceServer) mustEmbedUnimplementedLifeLogServiceServer() {}
func (UnimplementedLifeLogServiceServer) testEmbeddedByValue()                        {}

// UnsafeLifeLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LifeLogServiceServer will
// result in compilation errors.
type UnsafeLifeLogServiceServer interface {
	mustEmbedUnimplementedLifeLogServiceServer()
}

func RegisterLifeLogServiceServer(s grpc.ServiceRegistrar, srv LifeLogServiceServer) {
	// If the following call pancis, it indicates UnimplementedLifeLogServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LifeLogService_ServiceDesc, srv)
}

func _LifeLogService_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_Edit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).Edit(ctx, req.(*EditRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_SearchByTitle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchByTitleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).SearchByTitle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_SearchByTitle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).SearchByTitle(ctx, req.(*SearchByTitleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_DraftList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DraftListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).DraftList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_DraftList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).DraftList(ctx, req.(*DraftListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_Revoke_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).Revoke(ctx, req.(*RevokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LifeLogService_Detail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LifeLogServiceServer).Detail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LifeLogService_Detail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LifeLogServiceServer).Detail(ctx, req.(*DetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LifeLogService_ServiceDesc is the grpc.ServiceDesc for LifeLogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LifeLogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lifelog.v1.LifeLogService",
	HandlerType: (*LifeLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Edit",
			Handler:    _LifeLogService_Edit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _LifeLogService_Delete_Handler,
		},
		{
			MethodName: "SearchByTitle",
			Handler:    _LifeLogService_SearchByTitle_Handler,
		},
		{
			MethodName: "DraftList",
			Handler:    _LifeLogService_DraftList_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _LifeLogService_Revoke_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _LifeLogService_Publish_Handler,
		},
		{
			MethodName: "Detail",
			Handler:    _LifeLogService_Detail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/lifelog/v1/lifelog.proto",
}