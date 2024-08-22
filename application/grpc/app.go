package grpc

import (
	"antrein/bc-queue/application/common/repository"
	"antrein/bc-queue/internal/handler/analytic"
	"antrein/bc-queue/model/config"
	"context"

	pb "github.com/antrein/proto-repository/pb/bc"
	"google.golang.org/grpc"
)

type helloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func ApplicationDelegate(cfg *config.Config, repo *repository.CommonRepository) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	// Hello service
	helloServer := &helloServer{}
	pb.RegisterGreeterServer(grpcServer, helloServer)

	// Analytic service
	analyticServer := analytic.New(repo)
	pb.RegisterAnalyticServiceServer(grpcServer, analyticServer)

	return grpcServer, nil
}
