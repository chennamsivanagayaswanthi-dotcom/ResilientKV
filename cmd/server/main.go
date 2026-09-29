package main

import (
	"log"
	"net"

	rpcserver "resilientkv/grpc"
	pb "resilientkv/grpc/proto"
	"resilientkv/storage"

	"google.golang.org/grpc"
)

func main() {
	engine, err := storage.New("./data")
	if err != nil {
		log.Fatalf("failed to create storage engine: %v", err)
	}
	defer engine.Close()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	server := rpcserver.NewServer(engine)

	pb.RegisterKeyValueServiceServer(grpcServer, server)

	log.Println("ResilientKV gRPC server running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
