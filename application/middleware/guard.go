package guard

import (
	"antrein/bc-queue/model/config"
	"antrein/bc-queue/model/dto"
	"antrein/bc-queue/model/entity"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type GuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type AuthGuardContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Claims         entity.JWTClaim
}

func (g *GuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.NoBodyDTOResponseWrapper{
		Status:  status,
		Message: message,
	})
}

func (g *GuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultDTOResponseWrapper{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnError(status int, message string) error {
	g.ResponseWriter.WriteHeader(status)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.NoBodyDTOResponseWrapper{
		Status:  status,
		Message: message,
	})
}

func (g *AuthGuardContext) ReturnSuccess(data interface{}) error {
	g.ResponseWriter.WriteHeader(http.StatusOK)
	return json.NewEncoder(g.ResponseWriter).Encode(dto.DefaultDTOResponseWrapper{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func (g *AuthGuardContext) ReturnEvent(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(g.ResponseWriter, "data: %s\n\n", jsonData)
	if err != nil {
		return err // Handle writing errors
	}

	if flusher, ok := g.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	} else {
		return fmt.Errorf("streaming unsupported")
	}

	return nil
}

func DefaultGuard(handlerFunc func(g *GuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guardCtx := GuardContext{
			ResponseWriter: w,
			Request:        r,
		}
		if err := handlerFunc(&guardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func AuthGuard(cfg *config.Config, handlerFunc func(g *AuthGuardContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secrets.WaitingRoomSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		projectID, ok := claims["project_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID, ok := claims["session_id"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authGuardCtx := AuthGuardContext{
			ResponseWriter: w,
			Request:        r,
			Claims: entity.JWTClaim{
				ProjectID: projectID,
				SessionID: sessionID,
			},
		}

		if err := handlerFunc(&authGuardCtx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
