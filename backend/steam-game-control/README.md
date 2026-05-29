# Steam Game Server Control API

A Go-based API for controlling Steam game servers via systemctl with JWT authentication.

## Features

- Control Steam game servers using systemctl services
- JWT-based authentication for multiple users
- Real-time log streaming from systemd journal
- RESTful API endpoints for game management
- Configuration management with Viper
- User management with LevelDB and bcrypt

## Architecture

```
steam-game-control/
├── cmd/steamctl/        # CLI interface
├── configs/             # Configuration files
├── internal/
│   ├── app/             # Application server
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP request handlers
│   ├── service/         # Business logic
│   └── utils/           # Utility functions
└── go.mod               # Go module dependencies
```

## API Endpoints

### Authentication
- `POST /auth/login` - User login with JWT generation
- `POST /auth/register` - User registration

### Game Server Control
- `GET /games` - List available Steam games (systemctl services)
- `GET /games/{game}` - Get game server status
- `POST /games/{game}/start` - Start a game server
- `POST /games/{game}/stop` - Stop a game server
- `POST /games/{game}/restart` - Restart a game server
- `GET /games/{game}/logs` - Stream real-time logs

## User Management

The application uses LevelDB for storing user information and bcrypt for password hashing:

- **Storage**: LevelDB database for user records
- **Security**: bcrypt password hashing
- **Structure**: Each user stored with username, hashed password, and permissions

## Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Configure settings in `configs/config.yaml`
4. Run: `go run main.go`

## Configuration

The application uses Viper for configuration management:
- Default config path: `~/.steamctl/config.yaml`
- Environment variables can override settings

## Dependencies

- Cobra (CLI framework)
- Viper (Configuration)
- Gorilla Mux (HTTP routing)
- JWT (Authentication)
- go-systemd (systemctl integration)
- LevelDB (User storage)
- bcrypt (Password hashing)

## License

MIT