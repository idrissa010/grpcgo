package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "sami/rt0805/server/device_grpc" // Importez le fichier .pb.go généré

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// MongoDB configuration
const (
	mongoURI       = "mongodb://root:root@db:27017/"
	databaseName   = "devices"
	collectionName = "device_data"
)

// MongoDB client instance
var mongoClient *mongo.Client

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

		// Insert device data into MongoDB
		err = insertDeviceData(deviceData)
		if err != nil {
			log.Printf("Failed to insert device data into MongoDB: %v", err)
		}
	}
}

// Function to insert device data into MongoDB
func insertDeviceData(deviceData *pb.DeviceData) error {
	collection := mongoClient.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), deviceData)
	return err
}

func main() {
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())
	mongoClient = client

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
