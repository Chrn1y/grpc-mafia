// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: mafia.proto

package mafia_proto

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

// AppClient is the client API for App service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppClient interface {
	Play(ctx context.Context, opts ...grpc.CallOption) (App_PlayClient, error)
}

type appClient struct {
	cc grpc.ClientConnInterface
}

func NewAppClient(cc grpc.ClientConnInterface) AppClient {
	return &appClient{cc}
}

func (c *appClient) Play(ctx context.Context, opts ...grpc.CallOption) (App_PlayClient, error) {
	stream, err := c.cc.NewStream(ctx, &App_ServiceDesc.Streams[0], "/app.App/Play", opts...)
	if err != nil {
		return nil, err
	}
	x := &appPlayClient{stream}
	return x, nil
}

type App_PlayClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type appPlayClient struct {
	grpc.ClientStream
}

func (x *appPlayClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *appPlayClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AppServer is the server API for App service.
// All implementations must embed UnimplementedAppServer
// for forward compatibility
type AppServer interface {
	Play(App_PlayServer) error
	mustEmbedUnimplementedAppServer()
}

// UnimplementedAppServer must be embedded to have forward compatible implementations.
type UnimplementedAppServer struct {
}

func (UnimplementedAppServer) Play(App_PlayServer) error {
	return status.Errorf(codes.Unimplemented, "method Play not implemented")
}
func (UnimplementedAppServer) mustEmbedUnimplementedAppServer() {}

// UnsafeAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppServer will
// result in compilation errors.
type UnsafeAppServer interface {
	mustEmbedUnimplementedAppServer()
}

func RegisterAppServer(s grpc.ServiceRegistrar, srv AppServer) {
	s.RegisterService(&App_ServiceDesc, srv)
}

func _App_Play_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AppServer).Play(&appPlayServer{stream})
}

type App_PlayServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type appPlayServer struct {
	grpc.ServerStream
}

func (x *appPlayServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *appPlayServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// App_ServiceDesc is the grpc.ServiceDesc for App service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var App_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.App",
	HandlerType: (*AppServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Play",
			Handler:       _App_Play_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "mafia.proto",
}
