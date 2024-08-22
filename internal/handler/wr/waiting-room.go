package wr

import (
	"antrein/bc-queue/application/common/repository"
	guard "antrein/bc-queue/application/middleware"
	"antrein/bc-queue/internal/utils"
	"antrein/bc-queue/model/config"
	"antrein/bc-queue/model/dto"
	"antrein/bc-queue/model/entity"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Handler struct {
	cfg  *config.Config
	repo *repository.CommonRepository
}

func New(cfg *config.Config, repo *repository.CommonRepository) *Handler {
	return &Handler{
		cfg:  cfg,
		repo: repo,
	}
}

func (h *Handler) RegisterHandler(app *http.ServeMux) {
	app.HandleFunc("/bc/queue/register", guard.DefaultGuard(h.RegisterQueue))
	app.HandleFunc("/bc/queue/wr", guard.AuthGuard(h.cfg, h.UserQueue))
}

func (h *Handler) RegisterQueue(g *guard.GuardContext) error {
	ctx := context.Background()
	projectID := g.Request.URL.Query().Get("project_id")
	if projectID == "" {
		return g.ReturnError(500, "Project ID is missing")
	}
	config, err := h.repo.ConfigRepo.GetProjectConfig(ctx, projectID)
	if err != nil {
		return g.ReturnError(500, err.Error())
	}
	currentUser, err := h.repo.RoomRepo.CountUserInRoom(ctx, projectID, "main")
	if err != nil {
		fmt.Println("Redis Count User Error: ", err)
		return g.ReturnError(500, err.Error())
	}
	sessionID := uuid.New()
	session := entity.Session{
		SessionID:  sessionID.String(),
		EnqueuedAt: time.Now(),
	}
	waitingRoomClaim := entity.JWTClaim{
		SessionID: sessionID.String(),
		ProjectID: projectID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    projectID,
			Subject:   "",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	waitingRoomToken, err := utils.GenerateJWTToken(h.cfg.Secrets.WaitingRoomSecret, waitingRoomClaim)
	if err != nil {
		fmt.Println("Generate JWT Error: ", err)
		return g.ReturnError(500, err.Error())
	}
	// if currentUser <= int64(config.Threshold) || time.Now().Before(config.QueueEnd.AsTime()) || time.Now().After(config.QueueEnd.AsTime()) {
	if currentUser < int64(config.Threshold) {
		err = h.repo.RoomRepo.AddUserToMainRoom(ctx, projectID, session, int(config.SessionTime))
		if err != nil {
			fmt.Println("Adding User to Main Room: ", err)
			return g.ReturnError(500, err.Error())
		}
		mainRoomClaim := entity.JWTClaim{
			SessionID: sessionID.String(),
			ProjectID: projectID,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    projectID,
				Subject:   "",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.SessionTime))),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		mainRoomToken, err := utils.GenerateJWTToken(h.cfg.Secrets.MainRoomSecret, mainRoomClaim)
		if err != nil {
			fmt.Println("Generating Main Room Token Error: ", err)
			return g.ReturnError(500, err.Error())
		}
		tokens := dto.RegisterQueueResponse{
			WaitingRoomToken: waitingRoomToken,
			MainRoomToken:    mainRoomToken,
		}
		return g.ReturnSuccess(tokens)
	}
	err = h.repo.RoomRepo.AddUserToWaitingRoom(ctx, projectID, session)
	if err != nil {
		fmt.Println("Adding User to Waiting Room: ", err)
		return g.ReturnError(500, err.Error())
	}
	tokens := dto.RegisterQueueResponse{
		WaitingRoomToken: waitingRoomToken,
		MainRoomToken:    "",
	}
	return g.ReturnSuccess(tokens)
}

func (h *Handler) UserQueue(g *guard.AuthGuardContext) error {
	ctx := context.Background()
	g.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	g.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	g.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	g.ResponseWriter.Header().Set("Cache-Control", "no-cache")
	g.ResponseWriter.Header().Set("Connection", "keep-alive")

	sessionID := g.Claims.SessionID

	if sessionID == "" {
		return g.ReturnError(400, "Tidak terdaftar di queue")
	}

	projectID := g.Claims.ProjectID

	if projectID == "" {
		return g.ReturnError(400, "URL tidak terdaftar")
	}

	config, err := h.repo.ConfigRepo.GetProjectConfig(ctx, projectID)
	if err != nil {
		return g.ReturnError(500, err.Error())
	}

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return g.ReturnSuccess("Queue selesai")
		default:
			session, idx, err := h.repo.RoomRepo.GetUserFromRoom(ctx, projectID, "waiting", sessionID)
			if err != nil {
				fmt.Println("session", err.Error())
				return g.ReturnError(500, "Gagal mendapatkan data")
			}
			currentUser, err := h.repo.RoomRepo.CountUserInRoom(ctx, projectID, "main")
			if err != nil {
				fmt.Println("currentUser", err.Error())
				return g.ReturnError(500, "Gagal mendapatkan data")
			}
			if idx == 0 && currentUser < int64(config.Threshold) {
				err = h.repo.RoomRepo.AddUserToMainRoom(ctx, projectID, *session, int(config.SessionTime))
				if err != nil {
					return g.ReturnError(500, err.Error())
				}
				mainRoomClaim := entity.JWTClaim{
					SessionID: sessionID,
					ProjectID: projectID,
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    projectID,
						Subject:   "",
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.SessionTime))),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
					},
				}
				mainRoomToken, err := utils.GenerateJWTToken(h.cfg.Secrets.MainRoomSecret, mainRoomClaim)
				if err != nil {
					return g.ReturnError(500, err.Error())
				}
				event := dto.QueueEvent{
					IsFinished:    true,
					MainRoomToken: mainRoomToken,
					QueueNumber:   1,
					TimeRemaining: 0,
				}
				err = g.ReturnEvent(event)
				if err != nil {
					return g.ReturnError(500, err.Error())
				}
				return g.ReturnError(400, "Token tidak valid")
			} else {
				event := dto.QueueEvent{
					IsFinished:    false,
					MainRoomToken: "",
					QueueNumber:   idx + 1,
					TimeRemaining: float64(int32(idx+1) * config.SessionTime),
				}
				err = g.ReturnEvent(event)
				if err != nil {
					return g.ReturnError(500, err.Error())
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
	return g.ReturnSuccess("Queue selesai")
}
