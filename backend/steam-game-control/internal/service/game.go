package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/sdjournal"
)

// GameBackend defines the interface for game server operations.
// Implementations can use systemd D-Bus, mock backends for testing, etc.
type GameBackend interface {
	ListGames() ([]string, error)
	StartGame(name string) error
	StopGame(name string) error
	RestartGame(name string) error
	GetGameStatus(name string) (string, error)
	StreamLogs(name string, callback func(string)) error
}

// dbusConn abstracts the D-Bus connection for testability
type dbusConn interface {
	ListUnits() ([]dbus.UnitStatus, error)
	StartUnit(name string, mode string, ch chan<- string) (int, error)
	StopUnit(name string, mode string, ch chan<- string) (int, error)
	RestartUnit(name string, mode string, ch chan<- string) (int, error)
}

// GameService handles game server operations via systemd D-Bus
type GameService struct {
	conn dbusConn
}

// NewGameService creates a new GameService connected to systemd D-Bus
func NewGameService() *GameService {
	conn, err := dbus.New()
	if err != nil {
		panic(err)
	}

	return &GameService{
		conn: conn,
	}
}

// NewGameServiceMock creates a new GameService with a mocked D-Bus connection
func NewGameServiceMock(mock *dbusMock) *GameService {
	return &GameService{
		conn: mock,
	}
}

// NewGameServiceWithInterface creates a new GameService with any dbusConn implementation
func NewGameServiceWithInterface(conn dbusConn) *GameService {
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

	unitName := fmt.Sprintf("steam-%s.service", gameName)
	journal.AddMatch(sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT + "=" + unitName)

	journal.SeekTail()

	for {
		status := journal.Wait(time.Second)

		switch status {
		case sdjournal.SD_JOURNAL_APPEND:
			for {
				n, err := journal.Next()
				if err != nil {
					return err
				}
				if n == 0 {
					break
				}

				entry, err := journal.GetEntry()
				if err != nil {
					continue
				}

				logLine := fmt.Sprintf("[%d] %s", entry.RealtimeTimestamp, entry.Fields["MESSAGE"])
				callback(logLine)
			}
		case sdjournal.SD_JOURNAL_NOP:
			continue
		case sdjournal.SD_JOURNAL_INVALIDATE:
			return fmt.Errorf("journal invalidated")
		default:
			if status < 0 {
				return fmt.Errorf("error in Wait: %d", status)
			}
		}
	}
}

// ListUnitsResponse represents a systemd unit for mock testing
type ListUnitsResponse struct {
	Name        string
	ActiveState string
}

// dbusMock is a mock implementation of dbusConn for testing GameService
type dbusMock struct {
	units      []dbus.UnitStatus
	startErr   error
	stopErr    error
	restartErr error
}

func (m *dbusMock) ListUnits() ([]dbus.UnitStatus, error) {
	return m.units, nil
}

func (m *dbusMock) StartUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, m.startErr
}

func (m *dbusMock) StopUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, m.stopErr
}

func (m *dbusMock) RestartUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, m.restartErr
}

// MockGameBackend is a test double for GameBackend
type MockGameBackend struct {
	ListGamesFunc    func() ([]string, error)
	StartGameFunc    func(name string) error
	StopGameFunc     func(name string) error
	RestartGameFunc  func(name string) error
	GetGameStatusFunc func(name string) (string, error)
	StreamLogsFunc   func(name string, callback func(string)) error
}

// ListGames implements GameBackend
func (m *MockGameBackend) ListGames() ([]string, error) {
	return m.ListGamesFunc()
}

// StartGame implements GameBackend
func (m *MockGameBackend) StartGame(name string) error {
	return m.StartGameFunc(name)
}

// StopGame implements GameBackend
func (m *MockGameBackend) StopGame(name string) error {
	return m.StopGameFunc(name)
}

// RestartGame implements GameBackend
func (m *MockGameBackend) RestartGame(name string) error {
	return m.RestartGameFunc(name)
}

// GetGameStatus implements GameBackend
func (m *MockGameBackend) GetGameStatus(name string) (string, error) {
	return m.GetGameStatusFunc(name)
}

// StreamLogs implements GameBackend
func (m *MockGameBackend) StreamLogs(name string, callback func(string)) error {
	return m.StreamLogsFunc(name, callback)
}
