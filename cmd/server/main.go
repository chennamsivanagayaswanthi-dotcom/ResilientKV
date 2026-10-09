package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

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

	// Optional one-time network election at startup.
	if strings.EqualFold(os.Getenv("RAFT_START_ELECTION"), "true") {
		peersValue := os.Getenv("RAFT_PEERS")
		if strings.TrimSpace(peersValue) == "" {
			log.Fatal("RAFT_START_ELECTION=true requires RAFT_PEERS")
		}

		var peers []string
		for _, address := range strings.Split(peersValue, ",") {
			address = strings.TrimSpace(address)
			if address == "" {
				log.Fatal("RAFT_PEERS contains an empty address")
			}
			peers = append(peers, address)
		}

		log.Printf("Node %s starting network election; peers=%v", nodeID, peers)

		elected, electionErr := raftNode.ConductNetworkElection(
			peers,
			2*time.Second,
		)
		if electionErr != nil {
			log.Printf("Network election failed: %v", electionErr)
		} else if elected {
			log.Printf("Node %s won the election for term %d", nodeID, raftNode.GetCurrentTerm())
		} else {
			log.Printf("Node %s did not obtain a majority; state=%s term=%d",
				nodeID, raftNode.State, raftNode.GetCurrentTerm())
		}
	}

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
