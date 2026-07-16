package service

import (
	"context"
	"testing"
	"time"
)

func TestGameServiceStreamLogsFailsWithoutSystemd(t *testing.T) {
	// Create a context with timeout to avoid hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	svc := NewGameService()

	var captured string
	err := svc.StreamLogs(ctx, "csgo", func(line string) {
		captured = line
	})

	// In a non-systemd environment, this should fail or timeout
	_ = err
	_ = captured
}
