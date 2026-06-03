package service

import (
	"testing"
)

func TestGameServiceStreamLogsFailsWithoutSystemd(t *testing.T) {
	// Without a running systemd journal, NewJournal will fail.
	// This tests the initial error path of StreamLogs.
	svc := NewGameService()
	
	var captured string
	err := svc.StreamLogs("csgo", func(line string) {
		captured = line
	})
	
	// In a non-systemd environment, this should fail
	if err == nil {
		t.Fatal("expected error when streaming logs without systemd journal")
	}
	
	// No logs should be captured on error
	if captured != "" {
		t.Errorf("expected no captured logs, got '%s'", captured)
	}
}
