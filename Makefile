.PHONY: build-backend build-frontend run-backend run-frontend clean all help

BUILD_DIR = backend/steam-game-control
FRONTEND_DIR = frontend/games-frontend

# Build backend
build-backend:
	@echo "Building backend..."
	@cd $(BUILD_DIR) && go build -o steamctl main.go
	@echo "Backend built successfully"

# Build frontend
build-frontend:
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && npm install
	@cd $(FRONTEND_DIR) && npm run build
	@echo "Frontend built successfully"

# Build both
build: build-backend build-frontend
	@echo "All builds completed"

# Run backend
run-backend:
	@echo "Starting backend..."
	@if [ -f $(BUILD_DIR)/backend.pid ]; then \
		PID=$$(cat $(BUILD_DIR)/backend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Backend already running (PID: $$PID)"; \
			exit 1; \
		else \
			echo "Stale PID file found, removing..."; \
			rm $(BUILD_DIR)/backend.pid; \
		fi; \
	fi
	@cd $(BUILD_DIR) && go run main.go > $(BUILD_DIR)/backend.log 2>&1 & echo $$! > $(BUILD_DIR)/backend.pid
	@echo "Backend started (PID: $$(cat $(BUILD_DIR)/backend.pid))"

# Stop backend
stop-backend:
	@if [ -f $(BUILD_DIR)/backend.pid ]; then \
		PID=$$(cat $(BUILD_DIR)/backend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Stopping backend (PID: $$PID)..."; \
			kill $$PID; \
			rm $(BUILD_DIR)/backend.pid; \
			echo "Backend stopped"; \
		else \
			echo "Backend not running, removing stale PID file"; \
			rm $(BUILD_DIR)/backend.pid; \
		fi; \
	else \
		echo "No PID file found, backend not running"; \
	fi

# Run frontend
run-frontend:
	@echo "Starting frontend..."
	@if [ -f $(FRONTEND_DIR)/frontend.pid ]; then \
		PID=$$(cat $(FRONTEND_DIR)/frontend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Frontend already running (PID: $$PID)"; \
			exit 1; \
		else \
			echo "Stale PID file found, removing..."; \
			rm $(FRONTEND_DIR)/frontend.pid; \
		fi; \
	fi
	@cd $(FRONTEND_DIR) && npm run dev -- --host 0.0.0.0 > $(FRONTEND_DIR)/frontend.log 2>&1 & echo $$! > $(FRONTEND_DIR)/frontend.pid
	@sleep 2 && echo "Frontend started (PID: $$(cat $(FRONTEND_DIR)/frontend.pid))"

# Stop frontend
stop-frontend:
	@if [ -f $(FRONTEND_DIR)/frontend.pid ]; then \
		PID=$$(cat $(FRONTEND_DIR)/frontend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Stopping frontend (PID: $$PID)..."; \
			kill $$PID; \
			rm $(FRONTEND_DIR)/frontend.pid; \
			echo "Frontend stopped"; \
		else \
			echo "Frontend not running, removing stale PID file"; \
			rm $(FRONTEND_DIR)/frontend.pid; \
		fi; \
	else \
		echo "No PID file found, frontend not running"; \
	fi

# Run both backend and frontend
run: run-backend run-frontend
	@echo "Both services started"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5174"

# Stop both backend and frontend
stop: stop-backend stop-frontend
	@echo "All services stopped"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BUILD_DIR)/steamctl $(BUILD_DIR)/backend.pid $(BUILD_DIR)/backend.log
	@rm -f $(FRONTEND_DIR)/frontend.pid $(FRONTEND_DIR)/frontend.log
	@rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules $(FRONTEND_DIR)/.svelte-kit
	@echo "Cleaned"

# Show status
status:
	@echo "=== Backend Status ==="
	@if [ -f $(BUILD_DIR)/backend.pid ]; then \
		PID=$$(cat $(BUILD_DIR)/backend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Backend running (PID: $$PID)"; \
		else \
			echo "Backend not running (stale PID file)"; \
		fi; \
	else \
		echo "Backend not running"; \
	fi
	@echo ""
	@echo "=== Frontend Status ==="
	@if [ -f $(FRONTEND_DIR)/frontend.pid ]; then \
		PID=$$(cat $(FRONTEND_DIR)/frontend.pid); \
		if ps -p $$PID > /dev/null; then \
			echo "Frontend running (PID: $$PID)"; \
		else \
			echo "Frontend not running (stale PID file)"; \
		fi; \
	else \
		echo "Frontend not running"; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build both backend and frontend"
	@echo "  make build-backend - Build only backend"
	@echo "  make build-frontend - Build only frontend"
	@echo "  make run          - Start both backend and frontend"
	@echo "  make run-backend  - Start only backend"
	@echo "  make run-frontend - Start only frontend"
	@echo "  make stop         - Stop both backend and frontend"
	@echo "  make stop-backend - Stop only backend"
	@echo "  make stop-frontend - Stop only frontend"
	@echo "  make status       - Show status of both services"
	@echo "  make clean        - Clean build artifacts and PID files"
	@echo "  make help         - Show this help message"
