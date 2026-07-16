package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	cmd "dedctl/cmd/dedctl"
	"dedctl/internal/app"
	"dedctl/internal/config"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "hash" {
			// Let the hash command execute directly
			hashCmd := cmd.NewHashCmd()
			hashCmd.SetArgs(os.Args[2:])
			if err := hashCmd.Execute(); err != nil {
				os.Exit(1)
			}
			return
		}
		if arg == "server" || arg == "" {
			// Fall through to server startup
		} else {
			rootCmd.SetArgs(os.Args[1:])
			if err := rootCmd.Execute(); err != nil {
				os.Exit(1)
			}
			return
		}
	}

	fmt.Println("Initializing dedctl - Dedicated Game Controller...")

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--config" && i+1 < len(os.Args) {
			config.SetConfigFile(os.Args[i+1])
		}
	}

	err := app.Run()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("\nShutting down gracefully...")
}
