package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CommonResponse is the standard API response structure
type CommonResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest is the request body for register
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{
		"success": "false",
		"message": message,
	})
}

// CORS wraps a router with CORS middleware using the allowed origins set.
func CORS(router *mux.Router, allowedOrigins map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" || allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", r.Header.Get("Authorization"))
		router.ServeHTTP(w, r)
	})
}
