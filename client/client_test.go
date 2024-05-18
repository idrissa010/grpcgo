package main

import (
	"context"
	"testing"

	pb "sami/rt0805/client/device_grpc"

	"google.golang.org/grpc"
)

// Fonction pour rechercher les fichiers JSON dans un répertoire
var findJSONFiles = func(pattern string) ([]string, error) {
	return []string{pattern}, nil
}

// Fonction pour créer un client gRPC
var createDeviceServiceClient = func(cc grpc.ClientConnInterface) pb.DeviceServiceClient {
	return pb.NewDeviceServiceClient(cc)
}

type mockDeviceServiceClient struct{}

func (m *mockDeviceServiceClient) SendDeviceData(ctx context.Context, opts ...grpc.CallOption) (pb.DeviceService_SendDeviceDataClient, error) {
	return nil, nil
}

func TestMain(t *testing.T) {
	originalFindJSONFiles := findJSONFiles
	originalCreateDeviceServiceClient := createDeviceServiceClient

	findJSONFiles = func(pattern string) ([]string, error) {
		return []string{"donnees/test.json"}, nil
	}

	createDeviceServiceClient = func(cc grpc.ClientConnInterface) pb.DeviceServiceClient {
		return &mockDeviceServiceClient{}
	}
	defer func() {
		findJSONFiles = originalFindJSONFiles
		createDeviceServiceClient = originalCreateDeviceServiceClient
	}()

	main()

}
