package explore

import (
	"context"
	"log"

	"github.com/deepika-agrawal/explore/domain"
	"github.com/deepika-agrawal/explore/pb"
)

type ExploreServiceServer struct {
	pb.UnimplementedExploreServiceServer
	Repo domain.DecisionRepository
}

func (s *ExploreServiceServer) ListLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	log.Printf("Invoking ListLikedYou with req: %v", req)
	return s.Repo.ListLikedYou(ctx, req)
}

func (s *ExploreServiceServer) ListNewLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	log.Printf("Invoking ListNewLikedYou with req: %v", req)
	return s.Repo.ListNewLikedYou(ctx, req)
}

func (s *ExploreServiceServer) CountLikedYou(ctx context.Context, req *pb.CountLikedYouRequest) (*pb.CountLikedYouResponse, error) {
	log.Printf("Invoking CountLikedYou with req: %v", req)
	return s.Repo.CountLikedYou(ctx, req)
}

func (s *ExploreServiceServer) PutDecision(ctx context.Context, req *pb.PutDecisionRequest) (*pb.PutDecisionResponse, error) {
	log.Printf("Invoking PutDecision with req: %v", req)
	return s.Repo.PutDecision(ctx, req)
}
