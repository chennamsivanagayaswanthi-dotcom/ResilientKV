package grpc

import (
	"context"

	pb "resilientkv/grpc/proto"
	"resilientkv/storage"
)

type Server struct {
	pb.UnimplementedKeyValueServiceServer
	engine       *storage.Engine
	deduplicator *Deduplicator
}

func NewServer(engine *storage.Engine) *Server {
	return &Server{
		engine:       engine,
		deduplicator: NewDeduplicator(),
	}
}

func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	if req.GetRequestId() != "" && s.deduplicator.IsProcessed(req.GetRequestId()) {
		return &pb.PutResponse{
			Success: true,
		}, nil
	}

	if err := s.engine.Put(req.GetKey(), req.GetValue()); err != nil {
		return &pb.PutResponse{
			Success: false,
		}, err
	}

	if req.GetRequestId() != "" {
		s.deduplicator.MarkProcessed(req.GetRequestId())
	}

	return &pb.PutResponse{
		Success: true,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value, exists := s.engine.Get(req.GetKey())

	return &pb.GetResponse{
		Value: value,
		Found: exists,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.GetRequestId() != "" && s.deduplicator.IsProcessed(req.GetRequestId()) {
		return &pb.DeleteResponse{
			Success: true,
		}, nil
	}

	if err := s.engine.Delete(req.GetKey()); err != nil {
		return &pb.DeleteResponse{
			Success: false,
		}, err
	}

	if req.GetRequestId() != "" {
		s.deduplicator.MarkProcessed(req.GetRequestId())
	}

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}
