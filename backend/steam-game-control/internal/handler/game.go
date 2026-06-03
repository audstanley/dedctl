package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"steam-game-control/internal/service"
	"steam-game-control/internal/utils"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameService.ListGames()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list games: %v", err))
		return
	}

if games == nil {
		games = []string{}
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: "Games retrieved successfully",
		Data:    games,
	})
}

func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	err := h.gameService.StartGame(gameName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start game: %v", err))
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: fmt.Sprintf("Game %s started successfully", gameName),
		Data: map[string]string{
			"status": "started",
			"game":   gameName,
		},
	})
}

func (h *GameHandler) StopGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	err := h.gameService.StopGame(gameName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop game: %v", err))
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: fmt.Sprintf("Game %s stopped successfully", gameName),
		Data: map[string]string{
			"status": "stopped",
			"game":   gameName,
		},
	})
}

func (h *GameHandler) RestartGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	err := h.gameService.RestartGame(gameName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to restart game: %v", err))
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: fmt.Sprintf("Game %s restarted successfully", gameName),
		Data: map[string]string{
			"status": "restarted",
			"game":   gameName,
		},
	})
}

func (h *GameHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	sendLog := func(logLine string) {
		fmt.Fprintf(w, "data: %s\n\n", logLine)
		flusher.Flush()
	}

	err := h.gameService.StreamLogs(gameName, sendLog)
	if err != nil {
		fmt.Fprintf(w, "data: %s\n\n", err.Error())
		flusher.Flush()
		return
	}
}

func (h *GameHandler) GetGameStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	status, err := h.gameService.GetGameStatus(gameName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get game status: %v", err))
		return
	}

	if status == "" {
		status = "not-found"
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: "Game status retrieved",
		Data: map[string]string{
			"status": status,
			"game":   gameName,
		},
	})
}

func AuthRequired(secretKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")

			if token == "" {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					WriteError(w, http.StatusUnauthorized, "Authorization header required")
					return
				}

				tokenParts := strings.Split(authHeader, " ")
				if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
					WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
					return
				}

				token = tokenParts[1]
			}

			_, err := utils.ValidateToken(token, secretKey)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
