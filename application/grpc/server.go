package grpc

import (
	"antrein/bc-queue/model/config"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func StartServer(cfg *config.Config, grpcServer *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPC.Port))
	if err != nil {
		return err
	}
	return grpcServer.Serve(lis)
}
