package service

import (
	"context"
	pb "vado-client/internal/pb/server"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerService struct {
	pb.UnimplementedServerServiceServer
}

func (s *ServerService) Ping(_ context.Context, req *emptypb.Empty) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Run: true}, nil
}
