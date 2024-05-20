package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "sami/rt0805/server/device_grpc"
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

type mockMongoCollection struct {
	mock.Mock
}

func (m *mockMongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

type mockMongoClient struct {
	mock.Mock
}

func (m *mockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	args := m.Called(name, opts)
	return args.Get(0).(*mongo.Database)
}

type mockMongoDatabase struct {
	mock.Mock
}

func (m *mockMongoDatabase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	args := m.Called(name, opts)
	return args.Get(0).(*mongo.Collection)
}

func TestSendDeviceData(t *testing.T) {
	mockClient := new(mockMongoClient)
	mockDB := new(mockMongoDatabase)
	mockCollection := new(mockMongoCollection)

	mockClient.On("Database", databaseName, mock.Anything).Return(mockDB)
	mockDB.On("Collection", collectionName, mock.Anything).Return(mockCollection)
	mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, nil)

	mongoClient = &mongo.Client{} // This is just to avoid nil pointer dereference

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
		DeviceName: "device1",
		Operations: []*pb.Operation{
			{
				Type:         pb.Operation_CREATE,
				HasSucceeded: true,
			},
		},
	}
	if err := stream.Send(deviceData); err != nil {
		t.Fatalf("Failed to send device data: %v", err)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to close stream: %v", err)
	}
	mockCollection.AssertCalled(t, "InsertOne", mock.Anything, deviceData)
}
