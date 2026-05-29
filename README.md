# Audstanley Games - Game Server Control Platform

A full-stack game server management platform with a Go backend and Svelte frontend for controlling Steam game servers.

## Architecture

```
audstanley-games/
├── backend/
│   └── steam-game-control/      # Go backend API
│       ├── internal/
│       │   ├── app/              # Server setup
│       │   ├── config/           # Configuration
│       │   ├── handler/          # HTTP handlers
│       │   ├── service/          # Business logic
│       │   └── utils/            # Utilities
│       ├── cmd/                  # CLI interface
│       ├── configs/              # Config files
│       ├── go.mod                # Dependencies
│       └── main.go
│
├── frontend/
│   └── games-frontend/           # SvelteKit frontend
│       ├── src/
│       │   ├── lib/              # Shared components
│       │   └── routes/           # Page routes
│       ├── svelte.config.js
│       ├── vite.config.js
│       └── package.json
│
└── README.md
```

## Backend (Go)

The backend provides REST API endpoints for game server control:

- **Authentication**: JWT-based login/register
- **Game Control**: Start/stop/restart Steam game servers
- **Log Streaming**: Real-time log viewing via Server-Sent Events
- **User Management**: Multiple user accounts with admin support

### API Endpoints

```
POST   /auth/login           # User authentication
POST   /auth/register        # User registration  
GET    /games                # List available games (auth required)
POST   /games/{game}/start   # Start game server (auth required)
POST   /games/{game}/stop    # Stop game server (auth required)
POST   /games/{game}/restart # Restart game server (auth required)
GET    /games/{game}/logs    # Stream logs via SSE (auth required)
```

### Running the Backend

```bash
cd backend/steam-game-control
go mod tidy
go run main.go
```

Server runs on `http://localhost:8080` by default.

## Frontend (Svelte)

The frontend provides a web dashboard for managing game servers:

- **Login/Register**: User authentication interface
- **Dashboard**: View all available game servers
- **Game Controls**: Start/stop/restart servers from UI
- **Log Viewer**: Real-time log streaming

### Running the Frontend

```bash
cd frontend/games-frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173` by default.

## Technologies

### Backend
- Go 1.26.1
- Gorilla Mux (HTTP routing)
- Cobra (CLI)
- Viper (Configuration)
- JWT (Authentication)
- go-systemd (systemctl integration)
- LevelDB (User storage)

### Frontend
- SvelteKit
- Tailwind CSS
- Svelte stores (state management)
- EventSource API (log streaming)

## Implementation Phases

✅ **Phase 1**: Git repository setup with backend moved to `backend/`
✅ **Phase 2**: SvelteKit frontend with Tailwind CSS
✅ **Phase 3**: API client and authentication stores
✅ **Phase 4**: Login and registration pages
✅ **Phase 5**: Dashboard layout with navigation
✅ **Phase 6**: Game list dashboard
✅ **Phase 7**: Game detail page with controls
✅ **Phase 8**: Real-time log viewer with SSE

## License

## License

MIT
