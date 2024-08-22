package analytic

import (
	"antrein/bc-queue/application/common/repository"
	"context"
	"time"

	pb "github.com/antrein/proto-repository/pb/bc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedAnalyticServiceServer
	repo *repository.CommonRepository
}

func New(repo *repository.CommonRepository) *Server {
	return &Server{
		repo: repo,
	}
}

func (s *Server) GetAnalyticData(ctx context.Context, in *pb.AnalyticRequest) (*pb.AnalyticData, error) {
	projectID := in.GetProjectId()
	totalInQueue, err := s.repo.RoomRepo.CountUserInRoom(ctx, projectID, "waiting")
	if err != nil {
		return nil, err
	}
	totalInMain, err := s.repo.RoomRepo.CountUserInRoom(ctx, projectID, "main")
	if err != nil {
		return nil, err
	}
	return &pb.AnalyticData{
		ProjectId:         projectID,
		TotalUsersInQueue: int32(totalInQueue),
		TotalUsersInRoom:  int32(totalInMain),
		TotalUsers:        int32(totalInQueue) + int32(totalInMain),
		Timestamp:         timestamppb.Now(),
	}, nil
}

func (s *Server) StreamRealtimeData(in *pb.AnalyticRequest, stream pb.AnalyticService_StreamRealtimeDataServer) error {
	projectID := in.GetProjectId()
	ctx := context.Background()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case t := <-ticker.C:
			totalInQueue, err := s.repo.RoomRepo.CountUserInRoom(ctx, projectID, "waiting")
			if err != nil {
				return stream.Context().Err()
			}
			totalInMain, err := s.repo.RoomRepo.CountUserInRoom(ctx, projectID, "main")
			if err != nil {
				return stream.Context().Err()
			}
			analyticData := &pb.AnalyticData{
				ProjectId:         projectID,
				TotalUsersInQueue: int32(totalInQueue),
				TotalUsersInRoom:  int32(totalInMain),
				TotalUsers:        int32(totalInQueue) + int32(totalInMain),
				Timestamp:         timestamppb.New(t),
			}
			if err := stream.Send(analyticData); err != nil {
				return err
			}
		}
	}
}
