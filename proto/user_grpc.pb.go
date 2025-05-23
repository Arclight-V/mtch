// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: user.proto

package pb

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
	UserInfo_Register_FullMethodName = "/userService.UserInfo/Register"
	UserInfo_Login_FullMethodName    = "/userService.UserInfo/Login"
)

// UserInfoClient is the client API for UserInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserInfoClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type userInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewUserInfoClient(cc grpc.ClientConnInterface) UserInfoClient {
	return &userInfoClient{cc}
}

func (c *userInfoClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, UserInfo_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userInfoClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, UserInfo_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserInfoServer is the server API for UserInfo service.
// All implementations must embed UnimplementedUserInfoServer
// for forward compatibility.
type UserInfoServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedUserInfoServer()
}

// UnimplementedUserInfoServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserInfoServer struct{}

func (UnimplementedUserInfoServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserInfoServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserInfoServer) mustEmbedUnimplementedUserInfoServer() {}
func (UnimplementedUserInfoServer) testEmbeddedByValue()                  {}

// UnsafeUserInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserInfoServer will
// result in compilation errors.
type UnsafeUserInfoServer interface {
	mustEmbedUnimplementedUserInfoServer()
}

func RegisterUserInfoServer(s grpc.ServiceRegistrar, srv UserInfoServer) {
	// If the following call pancis, it indicates UnimplementedUserInfoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserInfo_ServiceDesc, srv)
}

func _UserInfo_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserInfoServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserInfo_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserInfoServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserInfo_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserInfoServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserInfo_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserInfoServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserInfo_ServiceDesc is the grpc.ServiceDesc for UserInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "userService.UserInfo",
	HandlerType: (*UserInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _UserInfo_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _UserInfo_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
