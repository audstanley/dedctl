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
	gameBackend service.GameBackend
}

func NewGameHandler(gameBackend service.GameBackend) *GameHandler {
	return &GameHandler{
		gameBackend: gameBackend,
	}
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameBackend.ListGames()
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

	err := h.gameBackend.StartGame(gameName)
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

	err := h.gameBackend.StopGame(gameName)
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

	err := h.gameBackend.RestartGame(gameName)
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

func escapeSSEData(data string) string {
	data = strings.ReplaceAll(data, "\\", "\\\\")
	data = strings.ReplaceAll(data, "\n", "\\n")
	data = strings.ReplaceAll(data, "\r", "\\r")
	return data
}

func (h *GameHandler) StreamLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Streaming not supported")
		return
	}

	ctx := r.Context()

	sendLog := func(logLine string) {
		select {
		case <-ctx.Done():
			return
		default:
			escaped := escapeSSEData(logLine)
			fmt.Fprintf(w, "data: %s\n\n", escaped)
			flusher.Flush()
		}
	}

	err := h.gameBackend.StreamLogs(r.Context(), gameName, sendLog)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", escapeSSEData(err.Error()))
			flusher.Flush()
		}
	}
}

func (h *GameHandler) GetGameStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameName := vars["game"]

	status, err := h.gameBackend.GetGameStatus(gameName)
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
