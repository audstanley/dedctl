package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/sdjournal"
)

// GameService handles game server operations
type GameService struct {
	conn *dbus.Conn
}

// NewGameService creates a new GameService
func NewGameService() *GameService {
	conn, err := dbus.New()
	if err != nil {
		panic(err)
	}

	return &GameService{
		conn: conn,
	}
}

// ListGames returns all available Steam games
func (s *GameService) ListGames() ([]string, error) {
	units, err := s.conn.ListUnits()
	if err != nil {
		return nil, err
	}

	var games []string
	for _, unit := range units {
		if strings.HasPrefix(unit.Name, "steam-") && strings.HasSuffix(unit.Name, ".service") {
			// Extract game name from service name (e.g., steam-csgo.service -> csgo)
			gameName := strings.TrimPrefix(strings.TrimSuffix(unit.Name, ".service"), "steam-")
			games = append(games, gameName)
		}
	}

	return games, nil
}

// StartGame starts a Steam game server
func (s *GameService) StartGame(gameName string) error {
	unitName := fmt.Sprintf("steam-%s.service", gameName)
	_, err := s.conn.StartUnit(unitName, "replace", nil)
	return err
}

// StopGame stops a Steam game server
func (s *GameService) StopGame(gameName string) error {
	unitName := fmt.Sprintf("steam-%s.service", gameName)
	_, err := s.conn.StopUnit(unitName, "replace", nil)
	return err
}

// RestartGame restarts a Steam game server
func (s *GameService) RestartGame(gameName string) error {
	unitName := fmt.Sprintf("steam-%s.service", gameName)
	_, err := s.conn.RestartUnit(unitName, "replace", nil)
	return err
}

// GetGameStatus returns the status of a Steam game server
func (s *GameService) GetGameStatus(gameName string) (string, error) {
	unitName := fmt.Sprintf("steam-%s.service", gameName)
	units, err := s.conn.ListUnits()
	if err != nil {
		return "", err
	}

	for _, unit := range units {
		if unit.Name == unitName {
			return string(unit.ActiveState), nil
		}
	}

	return "not-found", nil
}

// StreamLogs streams real-time logs from the systemd journal
func (s *GameService) StreamLogs(gameName string, callback func(string)) error {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return err
	}
	defer journal.Close()

	// Add filter for the specific game service
	unitName := fmt.Sprintf("steam-%s.service", gameName)
	journal.AddMatch(sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT + "=" + unitName)

	// Seek to the end of the journal
	journal.SeekTail()

	// Watch for new entries
	for {
		// Wait for new entries
		status := journal.Wait(time.Second)

		switch status {
		case sdjournal.SD_JOURNAL_APPEND:
			// New entries are available
			for {
				n, err := journal.Next()
				if err != nil {
					return err
				}
				if n == 0 {
					break // No more entries
				}

				// Get the entry
				entry, err := journal.GetEntry()
				if err != nil {
					continue
				}

				// Format and send the log entry
				logLine := fmt.Sprintf("[%d] %s", entry.RealtimeTimestamp, entry.Fields["MESSAGE"])
				callback(logLine)
			}
		case sdjournal.SD_JOURNAL_NOP:
			// No changes in journal, continue waiting
			continue
		case sdjournal.SD_JOURNAL_INVALIDATE:
			// Journal was invalidated, we should probably re-open it
			return fmt.Errorf("journal invalidated")
		default:
			if status < 0 {
				return fmt.Errorf("error in Wait: %d", status)
			}
		}
	}
}
