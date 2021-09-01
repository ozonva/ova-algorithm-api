// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ova_algorithm_api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OvaAlgorithmApiClient is the client API for OvaAlgorithmApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OvaAlgorithmApiClient interface {
	CreateAlgorithmV1(ctx context.Context, in *CreateAlgorithmRequestV1, opts ...grpc.CallOption) (*CreateAlgorithmResponseV1, error)
	DescribeAlgorithmV1(ctx context.Context, in *DescribeAlgorithmRequestV1, opts ...grpc.CallOption) (*DescribeAlgorithmResponseV1, error)
	ListAlgorithmsV1(ctx context.Context, in *ListAlgorithmsRequestV1, opts ...grpc.CallOption) (*ListAlgorithmsResponseV1, error)
	RemoveAlgorithmV1(ctx context.Context, in *RemoveAlgorithmRequestV1, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type ovaAlgorithmApiClient struct {
	cc grpc.ClientConnInterface
}

func NewOvaAlgorithmApiClient(cc grpc.ClientConnInterface) OvaAlgorithmApiClient {
	return &ovaAlgorithmApiClient{cc}
}

func (c *ovaAlgorithmApiClient) CreateAlgorithmV1(ctx context.Context, in *CreateAlgorithmRequestV1, opts ...grpc.CallOption) (*CreateAlgorithmResponseV1, error) {
	out := new(CreateAlgorithmResponseV1)
	err := c.cc.Invoke(ctx, "/ova.algorithm.api.OvaAlgorithmApi/CreateAlgorithmV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ovaAlgorithmApiClient) DescribeAlgorithmV1(ctx context.Context, in *DescribeAlgorithmRequestV1, opts ...grpc.CallOption) (*DescribeAlgorithmResponseV1, error) {
	out := new(DescribeAlgorithmResponseV1)
	err := c.cc.Invoke(ctx, "/ova.algorithm.api.OvaAlgorithmApi/DescribeAlgorithmV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ovaAlgorithmApiClient) ListAlgorithmsV1(ctx context.Context, in *ListAlgorithmsRequestV1, opts ...grpc.CallOption) (*ListAlgorithmsResponseV1, error) {
	out := new(ListAlgorithmsResponseV1)
	err := c.cc.Invoke(ctx, "/ova.algorithm.api.OvaAlgorithmApi/ListAlgorithmsV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ovaAlgorithmApiClient) RemoveAlgorithmV1(ctx context.Context, in *RemoveAlgorithmRequestV1, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/ova.algorithm.api.OvaAlgorithmApi/RemoveAlgorithmV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OvaAlgorithmApiServer is the server API for OvaAlgorithmApi service.
// All implementations must embed UnimplementedOvaAlgorithmApiServer
// for forward compatibility
type OvaAlgorithmApiServer interface {
	CreateAlgorithmV1(context.Context, *CreateAlgorithmRequestV1) (*CreateAlgorithmResponseV1, error)
	DescribeAlgorithmV1(context.Context, *DescribeAlgorithmRequestV1) (*DescribeAlgorithmResponseV1, error)
	ListAlgorithmsV1(context.Context, *ListAlgorithmsRequestV1) (*ListAlgorithmsResponseV1, error)
	RemoveAlgorithmV1(context.Context, *RemoveAlgorithmRequestV1) (*emptypb.Empty, error)
	mustEmbedUnimplementedOvaAlgorithmApiServer()
}

// UnimplementedOvaAlgorithmApiServer must be embedded to have forward compatible implementations.
type UnimplementedOvaAlgorithmApiServer struct {
}

func (UnimplementedOvaAlgorithmApiServer) CreateAlgorithmV1(context.Context, *CreateAlgorithmRequestV1) (*CreateAlgorithmResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlgorithmV1 not implemented")
}
func (UnimplementedOvaAlgorithmApiServer) DescribeAlgorithmV1(context.Context, *DescribeAlgorithmRequestV1) (*DescribeAlgorithmResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeAlgorithmV1 not implemented")
}
func (UnimplementedOvaAlgorithmApiServer) ListAlgorithmsV1(context.Context, *ListAlgorithmsRequestV1) (*ListAlgorithmsResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlgorithmsV1 not implemented")
}
func (UnimplementedOvaAlgorithmApiServer) RemoveAlgorithmV1(context.Context, *RemoveAlgorithmRequestV1) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAlgorithmV1 not implemented")
}
func (UnimplementedOvaAlgorithmApiServer) mustEmbedUnimplementedOvaAlgorithmApiServer() {}

// UnsafeOvaAlgorithmApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OvaAlgorithmApiServer will
// result in compilation errors.
type UnsafeOvaAlgorithmApiServer interface {
	mustEmbedUnimplementedOvaAlgorithmApiServer()
}

func RegisterOvaAlgorithmApiServer(s grpc.ServiceRegistrar, srv OvaAlgorithmApiServer) {
	s.RegisterService(&OvaAlgorithmApi_ServiceDesc, srv)
}

func _OvaAlgorithmApi_CreateAlgorithmV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlgorithmRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OvaAlgorithmApiServer).CreateAlgorithmV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ova.algorithm.api.OvaAlgorithmApi/CreateAlgorithmV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OvaAlgorithmApiServer).CreateAlgorithmV1(ctx, req.(*CreateAlgorithmRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _OvaAlgorithmApi_DescribeAlgorithmV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeAlgorithmRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OvaAlgorithmApiServer).DescribeAlgorithmV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ova.algorithm.api.OvaAlgorithmApi/DescribeAlgorithmV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OvaAlgorithmApiServer).DescribeAlgorithmV1(ctx, req.(*DescribeAlgorithmRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _OvaAlgorithmApi_ListAlgorithmsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAlgorithmsRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OvaAlgorithmApiServer).ListAlgorithmsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ova.algorithm.api.OvaAlgorithmApi/ListAlgorithmsV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OvaAlgorithmApiServer).ListAlgorithmsV1(ctx, req.(*ListAlgorithmsRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _OvaAlgorithmApi_RemoveAlgorithmV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveAlgorithmRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OvaAlgorithmApiServer).RemoveAlgorithmV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ova.algorithm.api.OvaAlgorithmApi/RemoveAlgorithmV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OvaAlgorithmApiServer).RemoveAlgorithmV1(ctx, req.(*RemoveAlgorithmRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

// OvaAlgorithmApi_ServiceDesc is the grpc.ServiceDesc for OvaAlgorithmApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OvaAlgorithmApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ova.algorithm.api.OvaAlgorithmApi",
	HandlerType: (*OvaAlgorithmApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAlgorithmV1",
			Handler:    _OvaAlgorithmApi_CreateAlgorithmV1_Handler,
		},
		{
			MethodName: "DescribeAlgorithmV1",
			Handler:    _OvaAlgorithmApi_DescribeAlgorithmV1_Handler,
		},
		{
			MethodName: "ListAlgorithmsV1",
			Handler:    _OvaAlgorithmApi_ListAlgorithmsV1_Handler,
		},
		{
			MethodName: "RemoveAlgorithmV1",
			Handler:    _OvaAlgorithmApi_RemoveAlgorithmV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ova-algorithm-api/ova-algorithm-api.proto",
}
