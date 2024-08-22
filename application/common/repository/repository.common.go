package repository

import (
	"antrein/bc-queue/application/common/resource"
	"antrein/bc-queue/internal/repository/config"
	"antrein/bc-queue/internal/repository/room"
	cfg "antrein/bc-queue/model/config"
)

type CommonRepository struct {
	ConfigRepo *config.Repository
	RoomRepo   *room.Repository
}

func NewCommonRepository(cfg *cfg.Config, rsc *resource.CommonResource) (*CommonRepository, error) {
	configRepo := config.New(cfg, rsc.Redis, rsc.GRPC)
	roomRepo := room.New(cfg, rsc.Redis)

	commonRepo := CommonRepository{
		ConfigRepo: configRepo,
		RoomRepo:   roomRepo,
	}
	return &commonRepo, nil
}
