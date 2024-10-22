// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: common/emote_rpc_service.proto

package commongrpc

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

// EmoteRpcServiceClient is the client API for EmoteRpcService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmoteRpcServiceClient interface {
	Subscribe(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (EmoteRpcService_SubscribeClient, error)
	SendSticker(ctx context.Context, in *StickerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendPing(ctx context.Context, in *EmotePingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type emoteRpcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmoteRpcServiceClient(cc grpc.ClientConnInterface) EmoteRpcServiceClient {
	return &emoteRpcServiceClient{cc}
}

func (c *emoteRpcServiceClient) Subscribe(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (EmoteRpcService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &EmoteRpcService_ServiceDesc.Streams[0], "/common.EmoteRpcService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &emoteRpcServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EmoteRpcService_SubscribeClient interface {
	Recv() (*EmoteEvent, error)
	grpc.ClientStream
}

type emoteRpcServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *emoteRpcServiceSubscribeClient) Recv() (*EmoteEvent, error) {
	m := new(EmoteEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *emoteRpcServiceClient) SendSticker(ctx context.Context, in *StickerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/common.EmoteRpcService/SendSticker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emoteRpcServiceClient) SendPing(ctx context.Context, in *EmotePingRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/common.EmoteRpcService/SendPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmoteRpcServiceServer is the server API for EmoteRpcService service.
// All implementations must embed UnimplementedEmoteRpcServiceServer
// for forward compatibility
type EmoteRpcServiceServer interface {
	Subscribe(*emptypb.Empty, EmoteRpcService_SubscribeServer) error
	SendSticker(context.Context, *StickerRequest) (*emptypb.Empty, error)
	SendPing(context.Context, *EmotePingRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedEmoteRpcServiceServer()
}

// UnimplementedEmoteRpcServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmoteRpcServiceServer struct {
}

func (UnimplementedEmoteRpcServiceServer) Subscribe(*emptypb.Empty, EmoteRpcService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedEmoteRpcServiceServer) SendSticker(context.Context, *StickerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSticker not implemented")
}
func (UnimplementedEmoteRpcServiceServer) SendPing(context.Context, *EmotePingRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPing not implemented")
}
func (UnimplementedEmoteRpcServiceServer) mustEmbedUnimplementedEmoteRpcServiceServer() {}

// UnsafeEmoteRpcServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmoteRpcServiceServer will
// result in compilation errors.
type UnsafeEmoteRpcServiceServer interface {
	mustEmbedUnimplementedEmoteRpcServiceServer()
}

func RegisterEmoteRpcServiceServer(s grpc.ServiceRegistrar, srv EmoteRpcServiceServer) {
	s.RegisterService(&EmoteRpcService_ServiceDesc, srv)
}

func _EmoteRpcService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EmoteRpcServiceServer).Subscribe(m, &emoteRpcServiceSubscribeServer{stream})
}

type EmoteRpcService_SubscribeServer interface {
	Send(*EmoteEvent) error
	grpc.ServerStream
}

type emoteRpcServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *emoteRpcServiceSubscribeServer) Send(m *EmoteEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _EmoteRpcService_SendSticker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StickerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmoteRpcServiceServer).SendSticker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.EmoteRpcService/SendSticker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmoteRpcServiceServer).SendSticker(ctx, req.(*StickerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmoteRpcService_SendPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmotePingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmoteRpcServiceServer).SendPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.EmoteRpcService/SendPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmoteRpcServiceServer).SendPing(ctx, req.(*EmotePingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmoteRpcService_ServiceDesc is the grpc.ServiceDesc for EmoteRpcService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmoteRpcService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "common.EmoteRpcService",
	HandlerType: (*EmoteRpcServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSticker",
			Handler:    _EmoteRpcService_SendSticker_Handler,
		},
		{
			MethodName: "SendPing",
			Handler:    _EmoteRpcService_SendPing_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _EmoteRpcService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "common/emote_rpc_service.proto",
}
