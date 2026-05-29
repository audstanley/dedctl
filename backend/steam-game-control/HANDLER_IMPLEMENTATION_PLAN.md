# HTTP Handler Implementation Plan

## Overview
Build the missing HTTP handler layer to expose the service layer as a RESTful API.

---

## Phase 1: Setup & Structure

**Create Directory:**
- `internal/handler/` directory

**Files to Create:**
- `internal/handler/common.go` - Common types, middleware, helpers
- `internal/handler/game.go` - Game control endpoints
- `internal/handler/auth.go` - Authentication endpoints

---

## Phase 2: Common Handler (internal/handler/common.go)

**Sudo Code:**
```
package handler

import (
    "encoding/json"
    "net/http"
    "steam-game-control/internal/utils"
)

// CommonResponse struct
- Success bool
- Message string
- Data interface{}

// WriteJSON helper function
- Takes http.ResponseWriter, statusCode, data
- Sets Content-Type header
- Encodes JSON response

// Middleware: AuthRequired
- Takes http.Handler
- Extracts Bearer token from Authorization header
- Validates token using utils.ValidateToken
- Returns 401 if invalid/missing
- Sets user context for valid tokens

// Request structures
- LoginRequest: Username, Password
- RegisterRequest: Username, Password, IsAdmin
- StartRequest: (optional) Config parameters
```

---

## Phase 3: Auth Handler (internal/handler/auth.go)

**Sudo Code:**
```
package handler

type AuthHandler struct {
    authService *service.AuthService
}

// NewAuthHandler constructor
- Takes authService
- Returns AuthHandler

// Login handler
- Parse JSON LoginRequest from body
- Call authService.Login(username, password)
- If error: return 400/401 with error message
- If success: return 200 with {token, user}
- Handle JSON decode errors

// Register handler
- Parse JSON RegisterRequest from body
- Call authService.Register(username, password, isAdmin)
- If error: return 400/409 (conflict) with error message
- If success: return 201 with success message
- Handle JSON decode errors

// GenerateJWT helper (if needed)
- Call utils.GenerateToken(username, secretKey)
- Return token string
```

---

## Phase 4: Game Handler (internal/handler/game.go)

**Sudo Code:**
```
package handler

type GameHandler struct {
    gameService *service.GameService
}

// NewGameHandler constructor
- Takes gameService
- Returns GameHandler

// ListGames handler
- AuthRequired middleware
- Call gameService.ListGames()
- If error: return 500 with error
- If success: return 200 with {games: []string}
- Handle empty list

// StartGame handler
- AuthRequired middleware
- Parse game name from URL params (mux.Vars)
- Call gameService.StartGame(gameName)
- If error: return 500 with error
- If success: return 200 with {status: "started", game: name}
- Extract vars from request

// StopGame handler
- AuthRequired middleware
- Parse game name from URL params
- Call gameService.StopGame(gameName)
- If error: return 500 with error
- If success: return 200 with {status: "stopped", game: name}

// RestartGame handler
- AuthRequired middleware
- Parse game name from URL params
- Call gameService.RestartGame(gameName)
- If error: return 500 with error
- If success: return 200 with {status: "restarted", game: name}

// StreamLogs handler
- AuthRequired middleware
- Parse game name from URL params
- Set Content-Type: text/event-stream for SSE
- Send headers for streaming
- Call gameService.StreamLogs(gameName, callback)
- Callback format: send JSON lines or raw log lines
- Handle client disconnect
```

---

## Phase 5: Error Handling & Response Format

**Sudo Code:**
```
// Error responses
- 400 Bad Request: Invalid input
- 401 Unauthorized: Missing/invalid token
- 403 Forbidden: Insufficient permissions
- 404 Not Found: Game not found
- 409 Conflict: User exists
- 500 Internal Server Error: Service failure

// Standard response format
{
    "success": true/false,
    "message": "description",
    "data": {...}
}
```

---

## Phase 6: Testing Checklist

- [ ] Auth: Register new user
- [ ] Auth: Login with valid credentials
- [ ] Auth: Login with invalid credentials (401)
- [ ] Auth: Protected endpoint without token (401)
- [ ] Games: List all games (empty list)
- [ ] Games: List all games (with games)
- [ ] Games: Start game (success)
- [ ] Games: Start game (game not found)
- [ ] Games: Stop game
- [ ] Games: Restart game
- [ ] Logs: Stream logs endpoint (connects)
- [ ] Logs: Stream logs (receives data)

---

## Dependencies to Verify

Check existing imports:
- `github.com/gorilla/mux` - URL params extraction
- `github.com/golang-jwt/jwt/v5` - JWT validation (already in utils)
- `steam-game-control/internal/service` - GameService, AuthService

---

## Notes

1. **StreamLogs implementation**: Consider using Server-Sent Events (SSE) for real-time streaming
2. **JWT secret**: Load from config in server.go, pass to handlers
3. **Logging**: Add request logging in middleware for debugging
4. **Graceful shutdown**: Ensure handlers handle context cancellation
5. **Password hashing**: Consider bcrypt for auth.go (currently plain text)

---

## Next Steps

1. Create `internal/handler/` directory
2. Implement `common.go` first (shared types/middleware)
3. Implement `auth.go` (simpler, no streaming)
4. Implement `game.go` (includes streaming)
5. Test endpoints manually with curl
6. Verify integration with server.go
