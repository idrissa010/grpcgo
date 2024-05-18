package main

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	pb "sami/rt0805/server/device_grpc"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterDeviceServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSendDeviceData(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mongoClient = mt.Client

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewDeviceServiceClient(conn)

	stream, err := client.SendDeviceData(context.Background())
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	deviceData := &pb.DeviceData{
		Id:   "test-id",
		Data: "test-data",
	}

	if err := stream.Send(deviceData); err != nil {
		t.Fatalf("Failed to send device data: %v", err)
	}

	if err := stream.CloseSend(); err != nil {
		t.Fatalf("Failed to close send: %v", err)
	}

	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to receive response: %v", err)
	}

	mt.AddMockResponses(mtest.CreateSuccessResponse())

	// Ensure the data was inserted
	var result pb.DeviceData
	err = mongoClient.Database(databaseName).Collection(collectionName).FindOne(context.Background(), map[string]interface{}{"id": deviceData.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find inserted data: %v", err)
	}

	if !proto.Equal(&result, deviceData) {
		t.Fatalf("Expected %v, got %v", deviceData, result)
	}
}

func TestInsertDeviceData(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mongoClient = mt.Client

	deviceData := &pb.DeviceData{
		Id:   "test-id",
		Data: "test-data",
	}

	err := insertDeviceData(deviceData)
	if err != nil {
		t.Fatalf("Failed to insert device data: %v", err)
	}

	mt.AddMockResponses(mtest.CreateSuccessResponse())

	// Ensure the data was inserted
	var result pb.DeviceData
	err = mongoClient.Database(databaseName).Collection(collectionName).FindOne(context.Background(), map[string]interface{}{"id": deviceData.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find inserted data: %v", err)
	}

	if !proto.Equal(&result, deviceData) {
		t.Fatalf("Expected %v, got %v", deviceData, result)
	}
}
