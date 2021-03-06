// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package auth

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

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Check(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Result, error)
	Create(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Result, error)
	Delete(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Nothing, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Check(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/auth.Auth/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Create(ctx context.Context, in *UserID, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/auth.Auth/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Delete(ctx context.Context, in *SessionID, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/auth.Auth/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations should embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	Check(context.Context, *SessionID) (*Result, error)
	Create(context.Context, *UserID) (*Result, error)
	Delete(context.Context, *SessionID) (*Nothing, error)
}

// UnimplementedAuthServer should be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) Check(context.Context, *SessionID) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAuthServer) Create(context.Context, *UserID) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAuthServer) Delete(context.Context, *SessionID) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.Auth/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Check(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.Auth/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Create(ctx, req.(*UserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.Auth/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Delete(ctx, req.(*SessionID))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Auth_Check_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Auth_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Auth_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "session.proto",
}
