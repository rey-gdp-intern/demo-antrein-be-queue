package resource

import (
	"antrein/bc-queue/model/config"
	"context"
	_ "database/sql"
	"log"

	"github.com/go-playground/validator"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CommonResource struct {
	Redis *redis.Client
	Vld   *validator.Validate
	GRPC  *grpc.ClientConn
}

func NewCommonResource(cfg *config.Config, ctx context.Context) (*CommonResource, error) {
	opt, err := redis.ParseURL(cfg.Database.RedisDB.URL)
	if err != nil {
		log.Println("Error parsing REDIS_URL:", err)
		return nil, err
	}
	redisClient := redis.NewClient(opt)
	err = redisClient.Ping(ctx).Err()

	if err != nil {
		log.Println("Error ping to redis:", err)
	}

	err = flushRedis(ctx, redisClient)
	if err != nil {
		log.Println("Error flushing redis:", err)
		return nil, err
	}

	vld := validator.New()

	grpcClient, err := grpc.Dial(cfg.GRPCConfig.DashboardQueue, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	rsc := CommonResource{
		Redis: redisClient,
		Vld:   vld,
		GRPC:  grpcClient,
	}
	return &rsc, nil
}

func flushRedis(ctx context.Context, redisClient *redis.Client) error {
	err := redisClient.FlushAll(ctx)
	return err.Err()
}
