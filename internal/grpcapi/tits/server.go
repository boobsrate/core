package titspbv1

import (
	"context"

	titsv1pb "github.com/boobsrate/apis/tits/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	titsv1pb.UnimplementedTitsServiceServer

	tits Service
}

func NewTitsGRPCServer(tits Service) *Server {
	return &Server{
		tits: tits,
	}
}

func (s Server) GetTits(ctx context.Context, request *titsv1pb.TitsRequest) (*titsv1pb.TitsResponse, error) {
	tits, err := s.tits.GetTits(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &titsv1pb.TitsResponse{
		Tits: titsListToProto(tits),
	}, nil
}

func (s Server) Vote(ctx context.Context, request *titsv1pb.VoteRequest) (*titsv1pb.TitsResponse, error) {
	err := s.tits.IncreaseRating(ctx, request.GetTitsId())
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	tits, err := s.tits.GetTits(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &titsv1pb.TitsResponse{
		Tits: titsListToProto(tits),
	}, nil

}

func (s *Server) Register(server *grpc.Server) {
	titsv1pb.RegisterTitsServiceServer(server, s)
}
