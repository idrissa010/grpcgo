# grpcgo


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=server/device_grpc --go-grpc_opt=paths=source_relative device.proto

go mod tidy