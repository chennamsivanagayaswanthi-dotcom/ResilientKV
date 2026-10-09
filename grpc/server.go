package grpc

import (
	"context"
	"sync"
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

	heartbeatMu     sync.Mutex
	heartbeatCancel context.CancelFunc
}

func NewServer(engine *storage.Engine) *Server {
	return newServer(engine, nil)
}

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

// TriggerElection starts an election and starts heartbeats if this node wins.
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

	if elected {
		s.heartbeatMu.Lock()

		if s.heartbeatCancel != nil {
			s.heartbeatCancel()
			s.heartbeatCancel = nil
		}

		heartbeatCtx, cancel := context.WithCancel(context.Background())
		stop, startErr := s.raftNode.StartHeartbeats(
			heartbeatCtx,
			req.GetPeerAddresses(),
			raft.HeartbeatConfig{
				Interval: 500 * time.Millisecond,
				Timeout:  300 * time.Millisecond,
			},
		)
		if startErr != nil {
			cancel()
			s.heartbeatMu.Unlock()
			return nil, status.Errorf(codes.Internal, "start heartbeats: %v", startErr)
		}

		s.heartbeatCancel = func() {
			cancel()
			stop()
		}
		s.heartbeatMu.Unlock()
	}

	return &pb.TriggerElectionResponse{
		Elected: elected,
		Term:    int32(s.raftNode.GetCurrentTerm()),
		State:   s.raftNode.GetState().String(),
	}, nil
}

// AppendEntries handles a Raft heartbeat or log replication request.
func (s *Server) AppendEntries(
	ctx context.Context,
	req *pb.AppendEntriesRequest,
) (*pb.AppendEntriesResponse, error) {
	if req == nil || req.GetLeaderId() == "" ||
		req.GetTerm() < 0 ||
		req.GetPrevLogIndex() < 0 ||
		req.GetPrevLogTerm() < 0 ||
		req.GetLeaderCommit() < 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid AppendEntries request")
	}

	if s.raftNode == nil {
		return nil, status.Error(codes.FailedPrecondition, "Raft node is not configured")
	}

	entries := make([]raft.LogEntry, 0, len(req.GetEntries()))
	for _, entry := range req.GetEntries() {
		if entry == nil || entry.GetTerm() < 0 || entry.GetTerm() > req.GetTerm() {
			return nil, status.Error(codes.InvalidArgument, "invalid replicated log entry")
		}

		entries = append(entries, raft.LogEntry{
			Term:    int(entry.GetTerm()),
			Command: entry.GetCommand(),
		})
	}

	result, err := s.raftNode.HandleAppendEntries(raft.AppendEntriesRequest{
		Term:         int(req.GetTerm()),
		LeaderID:     req.GetLeaderId(),
		PrevLogIndex: int(req.GetPrevLogIndex()),
		PrevLogTerm:  int(req.GetPrevLogTerm()),
		Entries:      entries,
		LeaderCommit: int(req.GetLeaderCommit()),
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "AppendEntries failed: %v", err)
	}

	return &pb.AppendEntriesResponse{
		Term:    int32(result.Term),
		Success: result.Success,
	}, nil
}
