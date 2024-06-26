// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: device.proto

package device_grpc

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

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	SendDeviceData(ctx context.Context, opts ...grpc.CallOption) (DeviceService_SendDeviceDataClient, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) SendDeviceData(ctx context.Context, opts ...grpc.CallOption) (DeviceService_SendDeviceDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeviceService_ServiceDesc.Streams[0], "/DeviceService/SendDeviceData", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceServiceSendDeviceDataClient{stream}
	return x, nil
}

type DeviceService_SendDeviceDataClient interface {
	Send(*DeviceData) error
	CloseAndRecv() (*SendResponse, error)
	grpc.ClientStream
}

type deviceServiceSendDeviceDataClient struct {
	grpc.ClientStream
}

func (x *deviceServiceSendDeviceDataClient) Send(m *DeviceData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *deviceServiceSendDeviceDataClient) CloseAndRecv() (*SendResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SendResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility
type DeviceServiceServer interface {
	SendDeviceData(DeviceService_SendDeviceDataServer) error
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServiceServer struct {
}

func (UnimplementedDeviceServiceServer) SendDeviceData(DeviceService_SendDeviceDataServer) error {
	return status.Errorf(codes.Unimplemented, "method SendDeviceData not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_SendDeviceData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DeviceServiceServer).SendDeviceData(&deviceServiceSendDeviceDataServer{stream})
}

type DeviceService_SendDeviceDataServer interface {
	SendAndClose(*SendResponse) error
	Recv() (*DeviceData, error)
	grpc.ServerStream
}

type deviceServiceSendDeviceDataServer struct {
	grpc.ServerStream
}

func (x *deviceServiceSendDeviceDataServer) SendAndClose(m *SendResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *deviceServiceSendDeviceDataServer) Recv() (*DeviceData, error) {
	m := new(DeviceData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendDeviceData",
			Handler:       _DeviceService_SendDeviceData_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "device.proto",
}
