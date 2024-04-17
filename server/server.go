package main

import (
	"fmt"
	"log"
	"net"

	pb "sami/rt0805/server/device_grpc" // Importez le fichier .pb.go généré

	"google.golang.org/grpc"
)

// server struct to implement DeviceServiceServer interface
type server struct {
	pb.UnimplementedDeviceServiceServer
}

// Implement the SendDeviceData method of the DeviceServiceServer interface
func (s *server) SendDeviceData(stream pb.DeviceService_SendDeviceDataServer) error {
	for {
		deviceData, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Received device data: %+v", deviceData)
	}
}

func main() {
	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Create a gRPC server instance
	s := grpc.NewServer()
	// Register the server instance with the DeviceServiceServer interface
	pb.RegisterDeviceServiceServer(s, &server{})
	// Start the gRPC server
	fmt.Println("Server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
