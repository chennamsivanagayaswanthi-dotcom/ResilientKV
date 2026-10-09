package main

import (
	"log"
	"net"
	"os"

	rpcserver "resilientkv/grpc"
	pb "resilientkv/grpc/proto"
	"resilientkv/raft"
	"resilientkv/storage"

	"google.golang.org/grpc"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	nodeID := envOrDefault("RAFT_NODE_ID", "node1")
	port := envOrDefault("RESILIENTKV_PORT", "50051")
	dataDir := envOrDefault("RESILIENTKV_DATA_DIR", "./data")

	engine, err := storage.New(dataDir)
	if err != nil {
		log.Fatalf("failed to create storage engine: %v", err)
	}
	defer engine.Close()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	raftNode := raft.NewNode(nodeID)
	grpcServer := grpc.NewServer()

	server := rpcserver.NewServerWithRaft(engine, raftNode)
	pb.RegisterKeyValueServiceServer(grpcServer, server)

	log.Printf(
		"ResilientKV server running: node=%s port=%s state=%s",
		nodeID,
		port,
		raftNode.State,
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
