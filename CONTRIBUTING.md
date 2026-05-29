# Contributing to Audstanley Games

Thank you for your interest in contributing! This document outlines the development setup and guidelines.

## Development Setup

### Prerequisites

- Go 1.26.1 or higher
- Node.js 18+ and npm
- Linux system with systemd (for backend)

### Backend Development

```bash
cd backend/steam-game-control
go mod tidy
go run main.go
```

### Frontend Development

```bash
cd frontend/games-frontend
npm install
npm run dev
```

## Project Structure

- `backend/` - Go backend API
- `frontend/` - SvelteKit frontend

## Development Workflow

1. Create a feature branch from `main`
2. Make your changes
3. Test locally
4. Submit a pull request

## Code Style

### Backend (Go)
- Follow Go formatting (`gofmt`)
- Add comments for exported functions
- Keep functions focused and single-purpose

### Frontend (Svelte)
- Use TypeScript where possible
- Follow Svelte component conventions
- Use Tailwind CSS for styling
- Keep components small and reusable

## Testing

Tests are currently manual. Add automated tests as the project grows.

## Commit Messages

Use clear, concise commit messages:
- `feat: add new login page`
- `fix: resolve auth token issue`
- `docs: update README`

## Questions?

Open an issue for questions or discussions.
