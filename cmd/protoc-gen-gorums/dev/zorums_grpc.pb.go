// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dev

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ZorumsServiceClient is the client API for ZorumsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ZorumsServiceClient interface {
	// GRPCCall plain gRPC call; testing that Gorums can ignore these, but that
	// they are added to the _grpc.pb.go generated file.
	GRPCCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	// QuorumCall plain.
	QuorumCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	// QuorumCall with per_node_arg option.
	QuorumCallPerNodeArg(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type zorumsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewZorumsServiceClient(cc grpc.ClientConnInterface) ZorumsServiceClient {
	return &zorumsServiceClient{cc}
}

func (c *zorumsServiceClient) GRPCCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dev.ZorumsService/GRPCCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zorumsServiceClient) QuorumCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dev.ZorumsService/QuorumCall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *zorumsServiceClient) QuorumCallPerNodeArg(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dev.ZorumsService/QuorumCallPerNodeArg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ZorumsServiceServer is the server API for ZorumsService service.
// All implementations must embed UnimplementedZorumsServiceServer
// for forward compatibility
type ZorumsServiceServer interface {
	// GRPCCall plain gRPC call; testing that Gorums can ignore these, but that
	// they are added to the _grpc.pb.go generated file.
	GRPCCall(context.Context, *Request) (*Response, error)
	// QuorumCall plain.
	QuorumCall(context.Context, *Request) (*Response, error)
	// QuorumCall with per_node_arg option.
	QuorumCallPerNodeArg(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedZorumsServiceServer()
}

// UnimplementedZorumsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedZorumsServiceServer struct {
}

func (*UnimplementedZorumsServiceServer) GRPCCall(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GRPCCall not implemented")
}
func (*UnimplementedZorumsServiceServer) QuorumCall(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuorumCall not implemented")
}
func (*UnimplementedZorumsServiceServer) QuorumCallPerNodeArg(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuorumCallPerNodeArg not implemented")
}
func (*UnimplementedZorumsServiceServer) mustEmbedUnimplementedZorumsServiceServer() {}

func RegisterZorumsServiceServer(s *grpc.Server, srv ZorumsServiceServer) {
	s.RegisterService(&_ZorumsService_serviceDesc, srv)
}

func _ZorumsService_GRPCCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZorumsServiceServer).GRPCCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ZorumsService/GRPCCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZorumsServiceServer).GRPCCall(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZorumsService_QuorumCall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZorumsServiceServer).QuorumCall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ZorumsService/QuorumCall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZorumsServiceServer).QuorumCall(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ZorumsService_QuorumCallPerNodeArg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ZorumsServiceServer).QuorumCallPerNodeArg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.ZorumsService/QuorumCallPerNodeArg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ZorumsServiceServer).QuorumCallPerNodeArg(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _ZorumsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dev.ZorumsService",
	HandlerType: (*ZorumsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GRPCCall",
			Handler:    _ZorumsService_GRPCCall_Handler,
		},
		{
			MethodName: "QuorumCall",
			Handler:    _ZorumsService_QuorumCall_Handler,
		},
		{
			MethodName: "QuorumCallPerNodeArg",
			Handler:    _ZorumsService_QuorumCallPerNodeArg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zorums.proto",
}
