package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"steam-game-control/internal/utils"
)

func authTestMiddleware(secret string, next http.Handler) http.Handler {
	return AuthRequired(secret, nil)(next)
}

func TestAuthRequiredNoHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/games", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	body := strings.TrimSpace(w.Body.String())
	if !strings.Contains(body, "Authorization header required") {
		t.Errorf("expected 'Authorization header required' in body, got: %s", body)
	}
}

func TestAuthRequiredMissingBearer(t *testing.T) {
	req := httptest.NewRequest("GET", "/games", nil)
	req.Header.Set("Authorization", "not-a-real-token")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredInvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/games", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredValidToken(t *testing.T) {
	token, _ := utils.GenerateToken("testuser", "test-secret")

	req := httptest.NewRequest("GET", "/games", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthRequiredExpiredToken(t *testing.T) {
	// Generate a valid token with one secret, validate with another
	token, _ := utils.GenerateToken("testuser", "wrong-secret")

	req := httptest.NewRequest("GET", "/games", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthRequiredCorrectUser(t *testing.T) {
	token, _ := utils.GenerateToken("admin", "test-secret")

	req := httptest.NewRequest("POST", "/games/csgo/start", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authTestMiddleware("test-secret", handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
