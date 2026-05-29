package handler

import (
	"encoding/json"
	"net/http"
	"steam-game-control/internal/config"
	"steam-game-control/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	config      *config.JWTConfig
}

func NewAuthHandler(authService *service.AuthService, cfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"username": user.Username,
				"is_admin": user.IsAdmin,
			},
		},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.authService.Register(req.Username, req.Password, req.IsAdmin)
	if err != nil {
		if err.Error() == "user already exists" {
			WriteError(w, http.StatusConflict, "User already exists")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	WriteJSON(w, http.StatusCreated, CommonResponse{
		Success: true,
		Message: "User registered successfully",
	})
}
