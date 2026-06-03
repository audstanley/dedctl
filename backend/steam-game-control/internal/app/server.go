package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"steam-game-control/internal/config"
	"steam-game-control/internal/handler"
	"steam-game-control/internal/service"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "steamctl",
		Short: "Steam Game Server Control API",
		Long:  `A CLI tool to control Steam game servers via systemctl`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Steam Game Server Control API")
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.steamctl/config.yaml)")

	return rootCmd
}

// Run starts the server
func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Initialize services
	gameService := service.NewGameService()

	var users []service.UserInfo
	for _, u := range cfg.Users {
		users = append(users, service.UserInfo{
			Username:     u.Username,
			PasswordHash: u.PasswordHash,
			PasswordType: u.PasswordType,
			IsAdmin:      u.IsAdmin,
		})
	}

	authService := service.NewAuthService(users, cfg.JWT.SecretKey)

	// Initialize handlers
	gameHandler := handler.NewGameHandler(gameService)
	authHandler := handler.NewAuthHandler(authService, &cfg.JWT)

	// Setup router
	r := mux.NewRouter()

	// Auth routes (no auth required)
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Game routes (with auth middleware)
	gameRouter := r.PathPrefix("/games").Subrouter()
	gameRouter.Use(handler.AuthRequired(cfg.JWT.SecretKey))
	gameRouter.HandleFunc("", gameHandler.ListGames).Methods("GET")
	gameRouter.HandleFunc("/{game}/start", gameHandler.StartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/stop", gameHandler.StopGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/restart", gameHandler.RestartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/logs", gameHandler.StreamLogs).Methods("GET")
	gameRouter.HandleFunc("/{game}/status", gameHandler.GetGameStatus).Methods("GET")

	// Build CORS allowed origins set
	allowedOrigins := make(map[string]bool)
	for _, origin := range cfg.Server.Origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins[origin] = true
		}
	}

	// CORS middleware wraps the entire router
	corsHandler := handler.CORS(r, allowedOrigins)

	// Start server
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: corsHandler,
	}

	// Create a channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting server on %s:%s\n", cfg.Server.Host, cfg.Server.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-done
	fmt.Println("\nShutting down server...")

	return nil
}
