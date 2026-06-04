package handler

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"steam-game-control/internal/config"
	"steam-game-control/internal/service"
	"steam-game-control/internal/utils"
)

type mockGameBackend struct {
	listGamesFunc           func() ([]string, error)
	listGamesWithMetaFunc   func() ([]service.GameInfo, error)
	startGameFunc           func(name string) error
	stopGameFunc            func(name string) error
	restartGameFunc         func(name string) error
	getGameStatusFunc       func(name string) (string, error)
	streamLogsFunc          func(ctx context.Context, name string, callback func(string)) error
	updateMetadataFunc      func(name string, appId, order int) error
	updateArtFunc           func(name string, appId int) error
}

func (m *mockGameBackend) ListGames() ([]string, error) {
	return m.listGamesFunc()
}

func (m *mockGameBackend) ListGamesWithMeta() ([]service.GameInfo, error) {
	if m.listGamesWithMetaFunc != nil {
		return m.listGamesWithMetaFunc()
	}
	// Default: return empty
	return []service.GameInfo{}, nil
}

func (m *mockGameBackend) StartGame(name string) error {
	return m.startGameFunc(name)
}

func (m *mockGameBackend) StopGame(name string) error {
	return m.stopGameFunc(name)
}

func (m *mockGameBackend) RestartGame(name string) error {
	return m.restartGameFunc(name)
}

func (m *mockGameBackend) GetGameStatus(name string) (string, error) {
	return m.getGameStatusFunc(name)
}

func (m *mockGameBackend) StreamLogs(ctx context.Context, name string, callback func(string)) error {
	return m.streamLogsFunc(ctx, name, callback)
}

func (m *mockGameBackend) UpdateMetadata(name string, appId, order int) error {
	if m.updateMetadataFunc != nil {
		return m.updateMetadataFunc(name, appId, order)
	}
	return nil
}

func (m *mockGameBackend) UpdateArt(name string, appId int) error {
	if m.updateArtFunc != nil {
		return m.updateArtFunc(name, appId)
	}
	return nil
}

func setupGameRouter(handler *GameHandler) *mux.Router {
	r := mux.NewRouter()
	gameRouter := r.PathPrefix("/games").Subrouter()
	gameRouter.HandleFunc("", handler.ListGames).Methods("GET")
	gameRouter.HandleFunc("/{game}/start", handler.StartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/stop", handler.StopGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/restart", handler.RestartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/logs", handler.StreamLogs).Methods("GET")
	gameRouter.HandleFunc("/{game}/status", handler.GetGameStatus).Methods("GET")
	return r
}

func TestListGamesSuccess(t *testing.T) {
	backend := &mockGameBackend{
		listGamesWithMetaFunc: func() ([]service.GameInfo, error) {
			return []service.GameInfo{
				{Name: "csgo", AppId: 730, Order: 1, HasImage: true},
				{Name: "rust", AppId: 252490, Order: 2, HasImage: false},
				{Name: "terraria", AppId: 105600, Order: 0, HasImage: false},
			}, nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected data to be []interface{}, got %T", resp.Data)
	}
	if len(data) != 3 {
		t.Errorf("expected 3 games, got %d", len(data))
	}
	first := data[0].(map[string]interface{})
	if first["name"] != "csgo" {
		t.Errorf("expected name 'csgo', got '%v'", first["name"])
	}
	if first["app_id"] != float64(730) {
		t.Errorf("expected app_id 730, got %v", first["app_id"])
	}
	if first["has_image"] != true {
		t.Error("expected has_image true")
	}
}

func TestListGamesError(t *testing.T) {
	backend := &mockGameBackend{
		listGamesWithMetaFunc: func() ([]service.GameInfo, error) {
			return nil, errors.New("dbus error")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestListGamesEmpty(t *testing.T) {
	backend := &mockGameBackend{
		listGamesWithMetaFunc: func() ([]service.GameInfo, error) {
			return []service.GameInfo{}, nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected data to be []interface{}, got %T", resp.Data)
	}
	if len(data) != 0 {
		t.Errorf("expected empty array, got %d items", len(data))
	}
}

func TestStartGameSuccess(t *testing.T) {
	backend := &mockGameBackend{
		startGameFunc: func(name string) error {
			if name != "csgo" {
				t.Errorf("expected game name 'csgo', got '%s'", name)
			}
			return nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/csgo/start", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Data == nil {
		t.Fatal("expected data in response")
	}
	data := func() map[string]string { m := resp.Data.(map[string]interface{}); result := make(map[string]string); for k, v := range m { result[k] = v.(string) }; return result }()
	if data["status"] != "started" {
		t.Errorf("expected status 'started', got '%s'", data["status"])
	}
}

func TestStartGameError(t *testing.T) {
	backend := &mockGameBackend{
		startGameFunc: func(name string) error {
			return errors.New("service not found")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/missing/start", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestStopGameSuccess(t *testing.T) {
	backend := &mockGameBackend{
		stopGameFunc: func(name string) error {
			return nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/rust/stop", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := func() map[string]string { m := resp.Data.(map[string]interface{}); result := make(map[string]string); for k, v := range m { result[k] = v.(string) }; return result }()
	if data["status"] != "stopped" {
		t.Errorf("expected status 'stopped', got '%s'", data["status"])
	}
}

func TestStopGameError(t *testing.T) {
	backend := &mockGameBackend{
		stopGameFunc: func(name string) error {
			return errors.New("permission denied")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/rust/stop", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRestartGameSuccess(t *testing.T) {
	backend := &mockGameBackend{
		restartGameFunc: func(name string) error {
			return nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/terraria/restart", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := func() map[string]string { m := resp.Data.(map[string]interface{}); result := make(map[string]string); for k, v := range m { result[k] = v.(string) }; return result }()
	if data["status"] != "restarted" {
		t.Errorf("expected status 'restarted', got '%s'", data["status"])
	}
}

func TestRestartGameError(t *testing.T) {
	backend := &mockGameBackend{
		restartGameFunc: func(name string) error {
			return errors.New("unit not found")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("POST", "/games/missing/restart", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetGameStatusSuccess(t *testing.T) {
	backend := &mockGameBackend{
		getGameStatusFunc: func(name string) (string, error) {
			return "active", nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games/csgo/status", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := func() map[string]string { m := resp.Data.(map[string]interface{}); result := make(map[string]string); for k, v := range m { result[k] = v.(string) }; return result }()
	if data["status"] != "active" {
		t.Errorf("expected status 'active', got '%s'", data["status"])
	}
}

func TestGetGameStatusNotFound(t *testing.T) {
	backend := &mockGameBackend{
		getGameStatusFunc: func(name string) (string, error) {
			return "", nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games/missing/status", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := func() map[string]string { m := resp.Data.(map[string]interface{}); result := make(map[string]string); for k, v := range m { result[k] = v.(string) }; return result }()
	if data["status"] != "not-found" {
		t.Errorf("expected status 'not-found', got '%s'", data["status"])
	}
}

func TestGetGameStatusError(t *testing.T) {
	backend := &mockGameBackend{
		getGameStatusFunc: func(name string) (string, error) {
			return "", errors.New("dbus error")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games/csgo/status", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestStreamLogsSuccess(t *testing.T) {
	backend := &mockGameBackend{
		streamLogsFunc: func(ctx context.Context, name string, callback func(string)) error {
			callback("[1234567890] Server started")
			callback("[1234567891] Map changed to de_dust2")
			return nil
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games/csgo/logs", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "text/event-stream" {
		t.Errorf("expected Content-Type 'text/event-stream', got '%s'", ct)
	}

	body := w.Body.String()
	if !strings.Contains(body, "data: [1234567890] Server started") {
		t.Errorf("expected log line in body, got: %s", body)
	}
	if !strings.Contains(body, "data: [1234567891] Map changed to de_dust2") {
		t.Errorf("expected second log line in body, got: %s", body)
	}
}

func TestStreamLogsError(t *testing.T) {
	backend := &mockGameBackend{
		streamLogsFunc: func(ctx context.Context, name string, callback func(string)) error {
			return errors.New("journal open failed")
		},
	}
	handler := NewGameHandler(backend)

	req := httptest.NewRequest("GET", "/games/missing/logs", nil)
	w := httptest.NewRecorder()
	router := setupGameRouter(handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "journal open failed") {
		t.Errorf("expected error message in body, got: %s", body)
	}
}

func TestAuthHandlerLoginSuccess(t *testing.T) {
	hash := fmt.Sprintf("%x", sha512.Sum512([]byte("admin123")))
	users := []service.UserInfo{
		{Username: "admin", PasswordHash: hash, PasswordType: "sha512", IsAdmin: true},
	}
	authService := service.NewAuthService(users, "test-secret")
	handler := NewAuthHandler(authService, &config.JWTConfig{SecretKey: "test-secret"})

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "admin123"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Data == nil {
		t.Fatal("expected data in response")
	}
	data := resp.Data.(map[string]interface{})
	if _, ok := data["token"]; !ok {
		t.Error("expected token in response data")
	}
	if user, ok := data["user"].(map[string]interface{}); ok {
		if user["username"] != "admin" {
			t.Errorf("expected username 'admin', got '%v'", user["username"])
		}
		if user["is_admin"] != true {
			t.Error("expected is_admin=true")
		}
	}
}

func TestAuthHandlerLoginInvalidBody(t *testing.T) {
	users := []service.UserInfo{
		{Username: "admin", PasswordHash: "hash", PasswordType: "sha512", IsAdmin: true},
	}
	authService := service.NewAuthService(users, "test-secret")
	handler := NewAuthHandler(authService, &config.JWTConfig{SecretKey: "test-secret"})

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAuthHandlerLoginWrongCredentials(t *testing.T) {
	users := []service.UserInfo{
		{Username: "admin", PasswordHash: "correct-hash", PasswordType: "sha512", IsAdmin: true},
	}
	authService := service.NewAuthService(users, "test-secret")
	handler := NewAuthHandler(authService, &config.JWTConfig{SecretKey: "test-secret"})

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrong"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestNewAuthHandler(t *testing.T) {
	users := []service.UserInfo{
		{Username: "admin", PasswordHash: "hash", PasswordType: "sha512", IsAdmin: true},
	}
	authService := service.NewAuthService(users, "secret")
	cfg := &config.JWTConfig{SecretKey: "secret", ExpiresIn: "24h"}

	h := NewAuthHandler(authService, cfg)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.authService == nil {
		t.Error("expected authService to be set")
	}
	if h.config == nil {
		t.Error("expected config to be set")
	}
}

func TestCORS(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	allowedOrigins := map[string]bool{
		"http://localhost:5174": true,
	}

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5174")
	w := httptest.NewRecorder()

	corsHandler := CORS(router, allowedOrigins)
	corsHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:5174" {
		t.Errorf("expected Access-Control-Allow-Origin 'http://localhost:5174', got '%s'", allowOrigin)
	}

	allowMethods := w.Header().Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		t.Error("expected Access-Control-Allow-Methods header")
	}

	allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
	if allowHeaders == "" {
		t.Error("expected Access-Control-Allow-Headers header")
	}

	allowCredentials := w.Header().Get("Access-Control-Allow-Credentials")
	if allowCredentials != "true" {
		t.Errorf("expected Access-Control-Allow-Credentials 'true', got '%s'", allowCredentials)
	}
}

func TestCORSAllowedOrigin(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	allowedOrigins := map[string]bool{
		"https://example.com": true,
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	corsHandler := CORS(router, allowedOrigins)
	corsHandler.ServeHTTP(w, req)

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "https://example.com" {
		t.Errorf("expected Access-Control-Allow-Origin 'https://example.com', got '%s'", allowOrigin)
	}
}

func TestCORSDisallowedOrigin(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	allowedOrigins := map[string]bool{
		"https://example.com": true,
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	corsHandler := CORS(router, allowedOrigins)
	corsHandler.ServeHTTP(w, req)

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "https://evil.com" {
		t.Error("should not set Access-Control-Allow-Origin for disallowed origin")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusOK, CommonResponse{
		Success: true,
		Message: "test",
		Data:    map[string]string{"key": "value"},
	})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	var resp CommonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, http.StatusBadRequest, "bad request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp["success"] != "false" {
		t.Error("expected success=false")
	}
	if resp["message"] != "bad request" {
		t.Errorf("expected message 'bad request', got '%s'", resp["message"])
	}
}

func TestAuthRequiredWithTokenQueryParam(t *testing.T) {
	token, _ := utils.GenerateToken("testuser", "test-secret")

	req := httptest.NewRequest("GET", "/games?token="+token, nil)
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AuthRequired("test-secret", nil)
	middleware(next).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with token query param, got %d", w.Code)
	}
}

func TestAuthRequiredTokenQueryParamInvalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/games?token=invalid-token", nil)
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AuthRequired("test-secret", nil)
	middleware(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with invalid token query param, got %d", w.Code)
	}
}
