// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pkg

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

// PuddleStoreClient is the client API for PuddleStore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PuddleStoreClient interface {
	ClientConnect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ClientID, error)
	ClientExit(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*Success, error)
	ClientOpen(ctx context.Context, in *OpenMessage, opts ...grpc.CallOption) (*OpenResponse, error)
	ClientClose(ctx context.Context, in *CloseMessage, opts ...grpc.CallOption) (*Success, error)
	ClientWrite(ctx context.Context, in *WriteMessage, opts ...grpc.CallOption) (*Success, error)
	ClientRead(ctx context.Context, in *ReadMessage, opts ...grpc.CallOption) (*ReadResponse, error)
	ClientMkdir(ctx context.Context, in *MkdirMessage, opts ...grpc.CallOption) (*Success, error)
	ClientRemove(ctx context.Context, in *RemoveMessage, opts ...grpc.CallOption) (*Success, error)
	ClientList(ctx context.Context, in *ListMessage, opts ...grpc.CallOption) (*ListResponse, error)
}

type puddleStoreClient struct {
	cc grpc.ClientConnInterface
}

func NewPuddleStoreClient(cc grpc.ClientConnInterface) PuddleStoreClient {
	return &puddleStoreClient{cc}
}

func (c *puddleStoreClient) ClientConnect(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ClientID, error) {
	out := new(ClientID)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientExit(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientExit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientOpen(ctx context.Context, in *OpenMessage, opts ...grpc.CallOption) (*OpenResponse, error) {
	out := new(OpenResponse)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientOpen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientClose(ctx context.Context, in *CloseMessage, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientClose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientWrite(ctx context.Context, in *WriteMessage, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientWrite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientRead(ctx context.Context, in *ReadMessage, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientRead", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientMkdir(ctx context.Context, in *MkdirMessage, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientMkdir", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientRemove(ctx context.Context, in *RemoveMessage, opts ...grpc.CallOption) (*Success, error) {
	out := new(Success)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *puddleStoreClient) ClientList(ctx context.Context, in *ListMessage, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/pkg.PuddleStore/ClientList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PuddleStoreServer is the server API for PuddleStore service.
// All implementations must embed UnimplementedPuddleStoreServer
// for forward compatibility
type PuddleStoreServer interface {
	ClientConnect(context.Context, *Empty) (*ClientID, error)
	ClientExit(context.Context, *ClientID) (*Success, error)
	ClientOpen(context.Context, *OpenMessage) (*OpenResponse, error)
	ClientClose(context.Context, *CloseMessage) (*Success, error)
	ClientWrite(context.Context, *WriteMessage) (*Success, error)
	ClientRead(context.Context, *ReadMessage) (*ReadResponse, error)
	ClientMkdir(context.Context, *MkdirMessage) (*Success, error)
	ClientRemove(context.Context, *RemoveMessage) (*Success, error)
	ClientList(context.Context, *ListMessage) (*ListResponse, error)
	mustEmbedUnimplementedPuddleStoreServer()
}

// UnimplementedPuddleStoreServer must be embedded to have forward compatible implementations.
type UnimplementedPuddleStoreServer struct {
}

func (UnimplementedPuddleStoreServer) ClientConnect(context.Context, *Empty) (*ClientID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientConnect not implemented")
}
func (UnimplementedPuddleStoreServer) ClientExit(context.Context, *ClientID) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientExit not implemented")
}
func (UnimplementedPuddleStoreServer) ClientOpen(context.Context, *OpenMessage) (*OpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientOpen not implemented")
}
func (UnimplementedPuddleStoreServer) ClientClose(context.Context, *CloseMessage) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientClose not implemented")
}
func (UnimplementedPuddleStoreServer) ClientWrite(context.Context, *WriteMessage) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientWrite not implemented")
}
func (UnimplementedPuddleStoreServer) ClientRead(context.Context, *ReadMessage) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientRead not implemented")
}
func (UnimplementedPuddleStoreServer) ClientMkdir(context.Context, *MkdirMessage) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientMkdir not implemented")
}
func (UnimplementedPuddleStoreServer) ClientRemove(context.Context, *RemoveMessage) (*Success, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientRemove not implemented")
}
func (UnimplementedPuddleStoreServer) ClientList(context.Context, *ListMessage) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientList not implemented")
}
func (UnimplementedPuddleStoreServer) mustEmbedUnimplementedPuddleStoreServer() {}

// UnsafePuddleStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PuddleStoreServer will
// result in compilation errors.
type UnsafePuddleStoreServer interface {
	mustEmbedUnimplementedPuddleStoreServer()
}

func RegisterPuddleStoreServer(s grpc.ServiceRegistrar, srv PuddleStoreServer) {
	s.RegisterService(&PuddleStore_ServiceDesc, srv)
}

func _PuddleStore_ClientConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientConnect(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientExit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientExit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientExit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientExit(ctx, req.(*ClientID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientOpen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientOpen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientOpen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientOpen(ctx, req.(*OpenMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientClose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientClose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientClose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientClose(ctx, req.(*CloseMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientWrite(ctx, req.(*WriteMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientRead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientRead(ctx, req.(*ReadMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientMkdir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MkdirMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientMkdir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientMkdir",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientMkdir(ctx, req.(*MkdirMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientRemove(ctx, req.(*RemoveMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PuddleStore_ClientList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PuddleStoreServer).ClientList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pkg.PuddleStore/ClientList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PuddleStoreServer).ClientList(ctx, req.(*ListMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// PuddleStore_ServiceDesc is the grpc.ServiceDesc for PuddleStore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PuddleStore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pkg.PuddleStore",
	HandlerType: (*PuddleStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClientConnect",
			Handler:    _PuddleStore_ClientConnect_Handler,
		},
		{
			MethodName: "ClientExit",
			Handler:    _PuddleStore_ClientExit_Handler,
		},
		{
			MethodName: "ClientOpen",
			Handler:    _PuddleStore_ClientOpen_Handler,
		},
		{
			MethodName: "ClientClose",
			Handler:    _PuddleStore_ClientClose_Handler,
		},
		{
			MethodName: "ClientWrite",
			Handler:    _PuddleStore_ClientWrite_Handler,
		},
		{
			MethodName: "ClientRead",
			Handler:    _PuddleStore_ClientRead_Handler,
		},
		{
			MethodName: "ClientMkdir",
			Handler:    _PuddleStore_ClientMkdir_Handler,
		},
		{
			MethodName: "ClientRemove",
			Handler:    _PuddleStore_ClientRemove_Handler,
		},
		{
			MethodName: "ClientList",
			Handler:    _PuddleStore_ClientList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/puddlestore.proto",
}
