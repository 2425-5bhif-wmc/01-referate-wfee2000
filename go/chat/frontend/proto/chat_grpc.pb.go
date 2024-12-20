// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.28.3
// source: chat.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ChatService_ClaimName_FullMethodName = "/chat.ChatService/ClaimName"
	ChatService_Connect_FullMethodName   = "/chat.ChatService/Connect"
)

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	ClaimName(ctx context.Context, in *ClaimNameRequest, opts ...grpc.CallOption) (*ClaimNameResponse, error)
	Connect(ctx context.Context, opts ...grpc.CallOption) (ChatService_ConnectClient, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) ClaimName(ctx context.Context, in *ClaimNameRequest, opts ...grpc.CallOption) (*ClaimNameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClaimNameResponse)
	err := c.cc.Invoke(ctx, ChatService_ClaimName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) Connect(ctx context.Context, opts ...grpc.CallOption) (ChatService_ConnectClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], ChatService_Connect_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceConnectClient{ClientStream: stream}
	return x, nil
}

type ChatService_ConnectClient interface {
	Send(*OutgoingMessage) error
	Recv() (*IncomingMessage, error)
	grpc.ClientStream
}

type chatServiceConnectClient struct {
	grpc.ClientStream
}

func (x *chatServiceConnectClient) Send(m *OutgoingMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceConnectClient) Recv() (*IncomingMessage, error) {
	m := new(IncomingMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	ClaimName(context.Context, *ClaimNameRequest) (*ClaimNameResponse, error)
	Connect(ChatService_ConnectServer) error
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) ClaimName(context.Context, *ClaimNameRequest) (*ClaimNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimName not implemented")
}
func (UnimplementedChatServiceServer) Connect(ChatService_ConnectServer) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_ClaimName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClaimNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).ClaimName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_ClaimName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).ClaimName(ctx, req.(*ClaimNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).Connect(&chatServiceConnectServer{ServerStream: stream})
}

type ChatService_ConnectServer interface {
	Send(*IncomingMessage) error
	Recv() (*OutgoingMessage, error)
	grpc.ServerStream
}

type chatServiceConnectServer struct {
	grpc.ServerStream
}

func (x *chatServiceConnectServer) Send(m *IncomingMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceConnectServer) Recv() (*OutgoingMessage, error) {
	m := new(OutgoingMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClaimName",
			Handler:    _ChatService_ClaimName_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _ChatService_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chat.proto",
}
