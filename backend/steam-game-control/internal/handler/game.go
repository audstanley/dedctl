package handler

import (
	"context"
	"encoding/json"
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
	infos, err := h.gameBackend.ListGamesWithMeta()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list games: %v", err))
		return
	}

	type gameInfoResponse struct {
		Name     string `json:"name"`
		AppId    int    `json:"app_id"`
		Order    int    `json:"order"`
		HasImage bool   `json:"has_image"`
	}

	respInfos := make([]gameInfoResponse, 0, len(infos))
	for _, info := range infos {
		respInfos = append(respInfos, gameInfoResponse{
			Name:     info.Name,
			AppId:    info.AppId,
			Order:    info.Order,
			HasImage: info.HasImage,
		})
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: "Games retrieved successfully",
		Data:    respInfos,
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

func (h *GameHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := r.Context().Value("is_admin").(bool)
	if !isAdmin {
		WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	gameName := vars["game"]

	var req struct {
		AppId int `json:"app_id"`
		Order int `json:"order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.gameBackend.UpdateMetadata(gameName, req.AppId, req.Order); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update metadata: %v", err))
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: fmt.Sprintf("Metadata for %s updated", gameName),
	})
}

func (h *GameHandler) UpdateArt(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := r.Context().Value("is_admin").(bool)
	if !isAdmin {
		WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	vars := mux.Vars(r)
	gameName := vars["game"]

	// Fetch the game's current metadata to get the AppID
	infos, err := h.gameBackend.ListGamesWithMeta()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list games: %v", err))
		return
	}

	var appId int
	for _, info := range infos {
		if info.Name == gameName {
			appId = info.AppId
			break
		}
	}

	if appId <= 0 {
		WriteError(w, http.StatusBadRequest, "No AppID set for this game")
		return
	}

	if err := h.gameBackend.UpdateArt(gameName, appId); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update game art: %v", err))
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: fmt.Sprintf("Game art for %s updated", gameName),
	})
}

func AuthRequired(secretKey string, users []service.UserInfo) mux.MiddlewareFunc {
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

			claims, err := utils.ValidateToken(token, secretKey)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			userInfo := resolveUserInfo(claims.Username, users)
			r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))
			r = r.WithContext(context.WithValue(r.Context(), "is_admin", userInfo.IsAdmin))

			next.ServeHTTP(w, r)
		})
	}
}

func resolveUserInfo(username string, users []service.UserInfo) *service.UserInfo {
	for _, u := range users {
		if u.Username == username {
			return &u
		}
	}
	return &service.UserInfo{Username: username}
}
