// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: route/groove.proto

package route

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	SymmetricKeyGen(ctx context.Context, in *SymmetricKeyGenRequest, opts ...grpc.CallOption) (*SymmetricKeyGenResponse, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) SymmetricKeyGen(ctx context.Context, in *SymmetricKeyGenRequest, opts ...grpc.CallOption) (*SymmetricKeyGenResponse, error) {
	out := new(SymmetricKeyGenResponse)
	err := c.cc.Invoke(ctx, "/Client/SymmetricKeyGen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServer is the server API for Client service.
// All implementations must embed UnimplementedClientServer
// for forward compatibility
type ClientServer interface {
	SymmetricKeyGen(context.Context, *SymmetricKeyGenRequest) (*SymmetricKeyGenResponse, error)
	mustEmbedUnimplementedClientServer()
}

// UnimplementedClientServer must be embedded to have forward compatible implementations.
type UnimplementedClientServer struct {
}

func (UnimplementedClientServer) SymmetricKeyGen(context.Context, *SymmetricKeyGenRequest) (*SymmetricKeyGenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SymmetricKeyGen not implemented")
}
func (UnimplementedClientServer) mustEmbedUnimplementedClientServer() {}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_SymmetricKeyGen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SymmetricKeyGenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).SymmetricKeyGen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Client/SymmetricKeyGen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).SymmetricKeyGen(ctx, req.(*SymmetricKeyGenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SymmetricKeyGen",
			Handler:    _Client_SymmetricKeyGen_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "route/groove.proto",
}

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerClient interface {
	CircuitSetup(ctx context.Context, in *CircuitSetupRequest, opts ...grpc.CallOption) (*CircuitSetupResponse, error)
	FetchPublicKey(ctx context.Context, in *FetchPublicKeyRequest, opts ...grpc.CallOption) (*FetchPublicKeyResponse, error)
}

type serverClient struct {
	cc grpc.ClientConnInterface
}

func NewServerClient(cc grpc.ClientConnInterface) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) CircuitSetup(ctx context.Context, in *CircuitSetupRequest, opts ...grpc.CallOption) (*CircuitSetupResponse, error) {
	out := new(CircuitSetupResponse)
	err := c.cc.Invoke(ctx, "/Server/CircuitSetup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) FetchPublicKey(ctx context.Context, in *FetchPublicKeyRequest, opts ...grpc.CallOption) (*FetchPublicKeyResponse, error) {
	out := new(FetchPublicKeyResponse)
	err := c.cc.Invoke(ctx, "/Server/FetchPublicKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
// All implementations must embed UnimplementedServerServer
// for forward compatibility
type ServerServer interface {
	CircuitSetup(context.Context, *CircuitSetupRequest) (*CircuitSetupResponse, error)
	FetchPublicKey(context.Context, *FetchPublicKeyRequest) (*FetchPublicKeyResponse, error)
	mustEmbedUnimplementedServerServer()
}

// UnimplementedServerServer must be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (UnimplementedServerServer) CircuitSetup(context.Context, *CircuitSetupRequest) (*CircuitSetupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CircuitSetup not implemented")
}
func (UnimplementedServerServer) FetchPublicKey(context.Context, *FetchPublicKeyRequest) (*FetchPublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchPublicKey not implemented")
}
func (UnimplementedServerServer) mustEmbedUnimplementedServerServer() {}

// UnsafeServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerServer will
// result in compilation errors.
type UnsafeServerServer interface {
	mustEmbedUnimplementedServerServer()
}

func RegisterServerServer(s grpc.ServiceRegistrar, srv ServerServer) {
	s.RegisterService(&Server_ServiceDesc, srv)
}

func _Server_CircuitSetup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CircuitSetupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).CircuitSetup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Server/CircuitSetup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).CircuitSetup(ctx, req.(*CircuitSetupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_FetchPublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchPublicKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).FetchPublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Server/FetchPublicKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).FetchPublicKey(ctx, req.(*FetchPublicKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Server_ServiceDesc is the grpc.ServiceDesc for Server service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Server_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CircuitSetup",
			Handler:    _Server_CircuitSetup_Handler,
		},
		{
			MethodName: "FetchPublicKey",
			Handler:    _Server_FetchPublicKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "route/groove.proto",
}
