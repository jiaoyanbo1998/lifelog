// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: api/proto/feed/feed.proto

package feedv1

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
	FeedService_CreateFeedEvent_FullMethodName = "/feed.v1.FeedService/CreateFeedEvent"
	FeedService_FindFeedEvents_FullMethodName  = "/feed.v1.FeedService/FindFeedEvents"
)

// FeedServiceClient is the client API for FeedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedServiceClient interface {
	CreateFeedEvent(ctx context.Context, in *CreateFeedEventRequest, opts ...grpc.CallOption) (*CreateFeedEventResponse, error)
	FindFeedEvents(ctx context.Context, in *FindFeedEventsRequest, opts ...grpc.CallOption) (*FindFeedEventsResponse, error)
}

type feedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedServiceClient(cc grpc.ClientConnInterface) FeedServiceClient {
	return &feedServiceClient{cc}
}

func (c *feedServiceClient) CreateFeedEvent(ctx context.Context, in *CreateFeedEventRequest, opts ...grpc.CallOption) (*CreateFeedEventResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFeedEventResponse)
	err := c.cc.Invoke(ctx, FeedService_CreateFeedEvent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) FindFeedEvents(ctx context.Context, in *FindFeedEventsRequest, opts ...grpc.CallOption) (*FindFeedEventsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FindFeedEventsResponse)
	err := c.cc.Invoke(ctx, FeedService_FindFeedEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedServiceServer is the server API for FeedService service.
// All implementations must embed UnimplementedFeedServiceServer
// for forward compatibility.
type FeedServiceServer interface {
	CreateFeedEvent(context.Context, *CreateFeedEventRequest) (*CreateFeedEventResponse, error)
	FindFeedEvents(context.Context, *FindFeedEventsRequest) (*FindFeedEventsResponse, error)
	mustEmbedUnimplementedFeedServiceServer()
}

// UnimplementedFeedServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFeedServiceServer struct{}

func (UnimplementedFeedServiceServer) CreateFeedEvent(context.Context, *CreateFeedEventRequest) (*CreateFeedEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeedEvent not implemented")
}
func (UnimplementedFeedServiceServer) FindFeedEvents(context.Context, *FindFeedEventsRequest) (*FindFeedEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindFeedEvents not implemented")
}
func (UnimplementedFeedServiceServer) mustEmbedUnimplementedFeedServiceServer() {}
func (UnimplementedFeedServiceServer) testEmbeddedByValue()                     {}

// UnsafeFeedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedServiceServer will
// result in compilation errors.
type UnsafeFeedServiceServer interface {
	mustEmbedUnimplementedFeedServiceServer()
}

func RegisterFeedServiceServer(s grpc.ServiceRegistrar, srv FeedServiceServer) {
	// If the following call pancis, it indicates UnimplementedFeedServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FeedService_ServiceDesc, srv)
}

func _FeedService_CreateFeedEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFeedEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).CreateFeedEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_CreateFeedEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).CreateFeedEvent(ctx, req.(*CreateFeedEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_FindFeedEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindFeedEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).FindFeedEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_FindFeedEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).FindFeedEvents(ctx, req.(*FindFeedEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedService_ServiceDesc is the grpc.ServiceDesc for FeedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "feed.v1.FeedService",
	HandlerType: (*FeedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFeedEvent",
			Handler:    _FeedService_CreateFeedEvent_Handler,
		},
		{
			MethodName: "FindFeedEvents",
			Handler:    _FeedService_FindFeedEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/feed/feed.proto",
}
