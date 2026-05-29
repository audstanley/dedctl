package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"steam-game-control/internal/app"
)

func main() {
	fmt.Println("Initializing Steam Game Server Control API...")

	// Create the server and run it
	err := app.Run()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

	// Wait for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("\nShutting down gracefully...")
}
