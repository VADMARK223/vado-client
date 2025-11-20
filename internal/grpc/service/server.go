package service

import (
	"context"
	pb "vado-client/api/pb/ping"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerService struct {
	pb.UnimplementedPingServiceServer
}

func (s *ServerService) Ping(_ context.Context, req *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Run: true}, nil
}
