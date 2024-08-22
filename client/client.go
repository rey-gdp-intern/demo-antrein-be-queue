package client

import (
	"context"
	"time"

	pb "github.com/antrein/proto-repository/pb/bc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Call(name string, grpcUrl string) (string, error) {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return r.GetMessage(), nil
}
