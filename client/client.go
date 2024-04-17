package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	pb "sami/rt0805/client/device_grpc" // Importez le fichier .pb.go généré

	"google.golang.org/grpc"
)

type Operation struct {
	Type         string `json:"type"`
	HasSucceeded bool   `json:"has_succeeded"`
}

type DeviceData struct {
	DeviceName   string      `json:"device_name"`
	Operations   []Operation `json:"operations"`
	NumOperation int32
	NumError     int32
}

func main() {
	// Connexion au serveur gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Impossible de se connecter : %v", err)
	}
	defer conn.Close()

	// Création du client gRPC
	client := pb.NewDeviceServiceClient(conn)

	// Liste des fichiers JSON dans le répertoire "donnees"
	files, err := filepath.Glob("donnees/*.json")
	if err != nil {
		fmt.Println("Erreur lors de la recherche des fichiers JSON :", err)
		return
	}

	// Parcourir chaque fichier JSON
	for _, file := range files {
		fmt.Println("Lecture du fichier :", file)

		// Lire le contenu du fichier
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier :", err)
			continue
		}

		// Parse du JSON en un tableau d'objets
		var deviceData []DeviceData
		err = json.Unmarshal(data, &deviceData)
		if err != nil {
			fmt.Println("Erreur lors du parsing du JSON :", err)
			continue
		}

		// Envoyer les données au serveur gRPC
		for _, device := range deviceData {
			// Créer un stream pour envoyer les données
			stream, err := client.SendDeviceData(context.Background())
			if err != nil {
				log.Fatalf("Impossible d'appeler SendDeviceData : %v", err)
			}

			// Convertir les opérations en []*device_grpc.Operation
			var operations []*pb.Operation
			for _, op := range device.Operations {
				var opType pb.Operation_OperationType
				switch op.Type {
				case "CREATE":
					opType = pb.Operation_CREATE
				case "UPDATE":
					opType = pb.Operation_UPDATE
				case "DELETE":
					opType = pb.Operation_DELETE
				default:
					log.Fatalf("Type d'opération invalide : %s", op.Type)
				}

				operations = append(operations, &pb.Operation{
					Type:         opType,
					HasSucceeded: op.HasSucceeded,
				})
			}

			// Envoyer chaque DeviceData au serveur
			err = stream.Send(&pb.DeviceData{
				DeviceName:    device.DeviceName,
				Operations:    operations, // Utiliser le type correct pour les opérations
				NumOperations: device.NumOperation,
				NumErrors:     device.NumError,
			})
			if err != nil {
				log.Fatalf("Impossible d'envoyer DeviceData : %v", err)
			}
		}
	}
}
