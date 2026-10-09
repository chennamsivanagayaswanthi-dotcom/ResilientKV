package grpc

import (
	"context"
	"time"

	pb "resilientkv/grpc/proto"
	"resilientkv/raft"
	"resilientkv/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedKeyValueServiceServer
	engine       *storage.Engine
	deduplicator *Deduplicator
	raftNode     *raft.Node
}

// NewServer preserves compatibility with the existing key-value server.
func NewServer(engine *storage.Engine) *Server {
	return newServer(engine, nil)
}

// NewServerWithRaft creates a server with Raft voting enabled.
func NewServerWithRaft(engine *storage.Engine, node *raft.Node) *Server {
	return newServer(engine, node)
}

func newServer(engine *storage.Engine, node *raft.Node) *Server {
	return &Server{
		engine:       engine,
		deduplicator: NewDeduplicator(),
		raftNode:     node,
	}
}

func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	if req.GetRequestId() != "" && s.deduplicator.IsProcessed(req.GetRequestId()) {
		return &pb.PutResponse{Success: true}, nil
	}

	if err := s.engine.Put(req.GetKey(), req.GetValue()); err != nil {
		return &pb.PutResponse{Success: false}, err
	}

	if req.GetRequestId() != "" {
		s.deduplicator.MarkProcessed(req.GetRequestId())
	}

	return &pb.PutResponse{Success: true}, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value, exists := s.engine.Get(req.GetKey())
	return &pb.GetResponse{Value: value, Found: exists}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.GetRequestId() != "" && s.deduplicator.IsProcessed(req.GetRequestId()) {
		return &pb.DeleteResponse{Success: true}, nil
	}

	if err := s.engine.Delete(req.GetKey()); err != nil {
		return &pb.DeleteResponse{Success: false}, err
	}

	if req.GetRequestId() != "" {
		s.deduplicator.MarkProcessed(req.GetRequestId())
	}

	return &pb.DeleteResponse{Success: true}, nil
}

// RequestVote handles an incoming Raft vote request.
func (s *Server) RequestVote(
	ctx context.Context,
	req *pb.RequestVoteRequest,
) (*pb.RequestVoteResponse, error) {
	if req == nil || req.GetCandidateId() == "" || req.GetTerm() < 0 ||
		req.GetLastLogIndex() < 0 || req.GetLastLogTerm() < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid vote request")
	}

	if s.raftNode == nil {
		return nil, status.Error(codes.FailedPrecondition, "Raft node is not configured")
	}

	granted := s.raftNode.RequestVote(
		req.GetCandidateId(),
		int(req.GetTerm()),
		int(req.GetLastLogIndex()),
		int(req.GetLastLogTerm()),
	)

	return &pb.RequestVoteResponse{
		Term:        int32(s.raftNode.GetCurrentTerm()),
		VoteGranted: granted,
	}, nil
}

// TriggerElection starts a network election on this running server's Raft node.
func (s *Server) TriggerElection(
	ctx context.Context,
	req *pb.TriggerElectionRequest,
) (*pb.TriggerElectionResponse, error) {
	if req == nil || len(req.GetPeerAddresses()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one peer address is required")
	}

	timeoutMillis := req.GetTimeoutMillis()
	if timeoutMillis <= 0 || timeoutMillis > 30000 {
		return nil, status.Error(codes.InvalidArgument, "timeout_millis must be between 1 and 30000")
	}

	if s.raftNode == nil {
		return nil, status.Error(codes.FailedPrecondition, "Raft node is not configured")
	}

	elected, err := s.raftNode.ConductNetworkElection(
		req.GetPeerAddresses(),
		time.Duration(timeoutMillis)*time.Millisecond,
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "election could not start: %v", err)
	}

	return &pb.TriggerElectionResponse{
		Elected: elected,
		Term:    int32(s.raftNode.GetCurrentTerm()),
		State:   s.raftNode.GetState().String(),
	}, nil
}
