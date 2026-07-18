# dedctl Backend

Go backend API and CLI tool for managing Steam dedicated game servers via systemd.

## Tech Stack

- **Go 1.26** — Core language
- **Cobra** (`github.com/spf13/cobra`) — CLI framework
- **Gorilla Mux** (`github.com/gorilla/mux`) — HTTP router
- **Viper** (`github.com/spf13/viper`) — Configuration management
- **go-systemd** (`github.com/coreos/go-systemd/v22`) — D-Bus and sdjournal integration
- **JWT** (`github.com/golang-jwt/jwt/v5`) — JSON Web Token authentication
- **bcrypt** (`golang.org/x/crypto`) — Password hashing

## Project Structure

```
backend/dedctl/
├── cmd/dedctl/           # CLI commands
│   ├── root.go           # Root command setup
│   ├── hash.go           # Password hash generation (sha512, bcrypt)
│   └── cache.go          # Bulk game cover image caching
├── internal/
│   ├── app/
│   │   └── server.go     # HTTP server setup, routing, middleware
│   ├── config/
│   │   ├── config.go     # YAML config loading (config.yaml)
│   │   ├── metadata.go   # Game metadata management (metadata.yaml)
│   │   └── *_test.go     # Config tests
│   ├── handler/
│   │   ├── common.go     # CORS middleware, common types
│   │   ├── auth.go       # Login handler, JWT auth middleware
│   │   ├── game.go       # Game CRUD, logs streaming, metadata updates
│   │   └── *_test.go     # Handler tests
│   ├── service/
│   │   ├── game.go       # Game server operations via systemd D-Bus
│   │   ├── auth.go       # Password verification and token generation
│   │   ├── images.go     # Game cover image download from Steam
│   │   └── *_test.go     # Service tests
│   └── utils/
│       ├── jwt.go        # JWT token generation/verification
│       └── jwt_test.go   # JWT tests
├── configs/              # Default configuration files
│   ├── config.yaml       # Server, JWT, game, and user config
│   ├── metadata.yaml     # Game metadata (AppIDs, display order, cover art)
│   └── img/              # Cached game cover images
├── main.go               # Entry point
├── go.mod
└── go.sum
```

## Configuration

### config.yaml

Located at `~/.dedctl/config.yaml` by default (or specified via `--config` flag):

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  origins:
    - "http://localhost:5174"

jwt:
  secret_key: "your-secret-key-here"
  expires_in: "24h"

game:
  base_path: "$HOME/Games"

users:
  - username: admin
    password_hash: "<hash>"
    password_type: sha512   # sha512, bcrypt, or plain
    is_admin: true
```

### metadata.yaml

Located at `~/.dedctl/metadata.yaml` (same directory as config.yaml), stores per-game metadata:

```yaml
main_image: main_cuttle.png
icon: cuttle_icon.png
games:
  <game_name>:
    app_id: 730
    order: 1
```

## CLI Commands

```bash
# Run the server
go run main.go server

# Generate a password hash
go run main.go hash <password> <type>
# Types: sha512 (default), bcrypt, plain

# Cache game images in bulk
go run main.go cache-images

# Custom config path
go run main.go --config /path/to/config.yaml server
```

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/login` | Authenticate and get JWT | No |
| GET | `/server-info` | Get global server metadata | No |
| GET | `/images/{name}` | Serve game cover images | No |
| GET | `/games` | List all game servers | Yes |
| POST | `/games/{game}/start` | Start a game server | Yes |
| POST | `/games/{game}/stop` | Stop a game server | Yes |
| POST | `/games/{game}/restart` | Restart a game server | Yes |
| GET | `/games/{game}/status` | Get game server status | Yes |
| GET | `/games/{game}/logs` | Stream logs (Server-Sent Events) | Yes |
| PATCH | `/games/{game}/metadata` | Update game AppID/order | Yes |
| POST | `/games/{game}/update-art` | Download game cover art | Yes |
| PATCH | `/games/settings` | Update global settings (main_image, icon) | Yes |

**Auth:** JWT Bearer token required. Include `Authorization: Bearer <token>` header.

## Running Tests

```bash
go test ./...
```

## Build

```bash
go build -o dedctl main.go
```

## Startup Flow

1. Load configuration from `config.yaml`
2. Ensure `metadata.yaml` and `img/` directory exist (auto-created if missing)
3. Load game metadata and scan systemd for `steam-*.service` units
4. Auto-add newly discovered games to metadata
5. Cache missing game cover images (by Steam AppID)
6. Initialize auth, game, and image services
7. Start HTTP server with CORS middleware on configured host:port
