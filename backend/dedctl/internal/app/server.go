package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"dedctl/internal/config"
	"dedctl/internal/handler"
	"dedctl/internal/service"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	configFile  string
	metadataDir string
)

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "dedctl",
		Short: "Dedctl - Dedicated Game Controller",
		Long:  `A CLI tool to manage dedicated Steam game servers`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dedctl")
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.dedctl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&metadataDir, "metadata", "", "metadata directory (default is same directory as config.yaml)")

	return rootCmd
}

// resolveMetadataDir returns the metadata directory path.
func resolveMetadataDir(cfgDir string) string {
	if metadataDir != "" {
		return metadataDir
	}
	return cfgDir
}

// Run starts the server
func Run() error {
	if configFile != "" {
		config.SetConfigFile(configFile)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Determine metadata directory
	var cfgDir string
	if cfgFileUsed := config.GetConfigFileUsed(); cfgFileUsed != "" {
		cfgDir = filepath.Dir(cfgFileUsed)
	} else {
		home, _ := os.UserHomeDir()
		cfgDir = filepath.Join(home, ".dedctl")
	}
	metaDir := resolveMetadataDir(cfgDir)
	imgDir := filepath.Join(metaDir, "img")

	if err := ensureDefaultImages(imgDir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to ensure default images: %v\n", err)
	}

	// Ensure metadata.yaml and img/ directory exist
	metaPath := filepath.Join(metaDir, "metadata.yaml")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		fmt.Printf("Creating metadata.yaml at %s\n", metaPath)
		if err := config.SaveMetadata(metaDir, &config.Metadata{Games: make(map[string]config.GameMetadata)}); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load metadata
	meta, err := config.LoadMetadata(metaDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		meta = &config.Metadata{Games: make(map[string]config.GameMetadata)}
	}

	// Initialize services
	gameService := service.NewGameService()
	imageService := service.NewImageService()

	// Set package-level metadata/image references for ListGamesWithMeta
	gameMetaMap := make(map[string]struct{ AppId int; Order int })
	for name, gm := range meta.Games {
		gameMetaMap[name] = struct{ AppId int; Order int }{AppId: gm.AppId, Order: gm.Order}
	}
	service.SetMetadataAndImages(gameMetaMap, imageService, imgDir)
	service.SetMetaDir(metaDir, meta)

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

	// Cache missing images for games that have app_id (no auto-add)
	games, err := gameService.ListGames()
	if err == nil && len(games) > 0 {
		fmt.Println("Checking for missing game images...")
		if err := imageService.CacheMissingImages(games, gameMetaMap, imgDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}

	// Initialize handlers
	gameHandler := handler.NewGameHandler(gameService)
	authHandler := handler.NewAuthHandler(authService, &cfg.JWT)

	// Setup router
	r := mux.NewRouter()

	// Auth routes (no auth required)
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Static image serving
	r.HandleFunc("/images/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		http.ServeFile(w, r, filepath.Join(imgDir, name))
	})

	r.HandleFunc("/server-info", gameHandler.GetServerInfo).Methods("GET")

	// Game routes (with auth middleware)
	gameRouter := r.PathPrefix("/games").Subrouter()
	gameRouter.Use(handler.AuthRequired(cfg.JWT.SecretKey, users))
	gameRouter.HandleFunc("", gameHandler.ListGames).Methods("GET")
	gameRouter.HandleFunc("/{game}/start", gameHandler.StartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/stop", gameHandler.StopGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/restart", gameHandler.RestartGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/enable", gameHandler.EnableGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/disable", gameHandler.DisableGame).Methods("POST")
	gameRouter.HandleFunc("/{game}/logs", gameHandler.StreamLogs).Methods("GET")
	gameRouter.HandleFunc("/{game}/status", gameHandler.GetGameStatus).Methods("GET")
	gameRouter.HandleFunc("/{game}/metadata", gameHandler.UpdateMetadata).Methods("PATCH")
	gameRouter.HandleFunc("/{game}/update-art", gameHandler.UpdateArt).Methods("POST")
	gameRouter.HandleFunc("/settings", gameHandler.UpdateGlobalSettings).Methods("PATCH")

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
