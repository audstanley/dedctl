package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
)

func TestGameServiceListGames(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "steam-csgo.service", ActiveState: "active"},
			{Name: "steam-rust.service", ActiveState: "inactive"},
			{Name: "some-other.service", ActiveState: "active"},
		},
	}
	svc := NewGameServiceMock(mock)

	games, err := svc.ListGames()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 2 {
		t.Errorf("expected 2 games, got %d", len(games))
	}
	if games[0] != "csgo" {
		t.Errorf("expected first game 'csgo', got '%s'", games[0])
	}
	if games[1] != "rust" {
		t.Errorf("expected second game 'rust', got '%s'", games[1])
	}
}

func TestGameServiceListGamesEmpty(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{},
	}
	svc := NewGameServiceMock(mock)

	games, err := svc.ListGames()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 0 {
		t.Errorf("expected 0 games, got %d", len(games))
	}
}

func TestGameServiceListGamesNoSteam(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "nginx.service", ActiveState: "active"},
			{Name: "docker.service", ActiveState: "active"},
		},
	}
	svc := NewGameServiceMock(mock)

	games, err := svc.ListGames()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 0 {
		t.Errorf("expected 0 games, got %d", len(games))
	}
}

func TestGameServiceStartGame(t *testing.T) {
	mock := &dbusMock{}
	svc := NewGameServiceMock(mock)

	err := svc.StartGame("csgo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGameServiceStartGameError(t *testing.T) {
	mock := &dbusMock{
		startErr: errors.New("unit not found"),
	}
	svc := NewGameServiceMock(mock)

	err := svc.StartGame("missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unit not found") {
		t.Errorf("expected 'unit not found' in error, got: %v", err)
	}
}

func TestGameServiceStopGame(t *testing.T) {
	mock := &dbusMock{}
	svc := NewGameServiceMock(mock)

	err := svc.StopGame("rust")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGameServiceStopGameError(t *testing.T) {
	mock := &dbusMock{
		stopErr: errors.New("permission denied"),
	}
	svc := NewGameServiceMock(mock)

	err := svc.StopGame("rust")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGameServiceRestartGame(t *testing.T) {
	mock := &dbusMock{}
	svc := NewGameServiceMock(mock)

	err := svc.RestartGame("terraria")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGameServiceRestartGameError(t *testing.T) {
	mock := &dbusMock{
		restartErr: errors.New("service not running"),
	}
	svc := NewGameServiceMock(mock)

	err := svc.RestartGame("terraria")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGameServiceGetGameStatusActive(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "steam-csgo.service", ActiveState: "active"},
		},
	}
	svc := NewGameServiceMock(mock)

	status, err := svc.GetGameStatus("csgo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != "active" {
		t.Errorf("expected status 'active', got '%s'", status)
	}
}

func TestGameServiceGetGameStatusInactive(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "steam-rust.service", ActiveState: "inactive"},
		},
	}
	svc := NewGameServiceMock(mock)

	status, err := svc.GetGameStatus("rust")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != "inactive" {
		t.Errorf("expected status 'inactive', got '%s'", status)
	}
}

func TestGameServiceGetGameStatusNotFound(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "steam-csgo.service", ActiveState: "active"},
		},
	}
	svc := NewGameServiceMock(mock)

	status, err := svc.GetGameStatus("missing")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != "not-found" {
		t.Errorf("expected status 'not-found', got '%s'", status)
	}
}

func TestGameServiceGetGameStatusError(t *testing.T) {
	mock2 := &dbusMockErr{
		listErr: errors.New("dbus connection failed"),
	}
	svc := NewGameServiceWithInterface(mock2)

	_, err := svc.GetGameStatus("csgo")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGameServiceNewGameServiceMock(t *testing.T) {
	mock := &dbusMock{}
	svc := NewGameServiceMock(mock)
	if svc == nil {
		t.Fatal("expected non-nil GameService")
	}
	if svc.conn != mock {
		t.Fatal("expected GameService to use mock connection")
	}
}

func TestGameServiceListGamesWithSuffix(t *testing.T) {
	mock := &dbusMock{
		units: []dbus.UnitStatus{
			{Name: "steam-csgo.service", ActiveState: "active"},
			{Name: "steam-counter-strike.service", ActiveState: "active"},
			{Name: "steam-dota2.service", ActiveState: "active"},
			{Name: "steam-csgo", ActiveState: "active"},
		},
	}
	svc := NewGameServiceMock(mock)

	games, err := svc.ListGames()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 3 {
		t.Errorf("expected 3 games, got %d", len(games))
	}
}

// dbusMockErr is a mock that returns errors on ListUnits
type dbusMockErr struct {
	listErr error
}

func (m *dbusMockErr) ListUnits() ([]dbus.UnitStatus, error) {
	return nil, m.listErr
}

func (m *dbusMockErr) StartUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, nil
}

func (m *dbusMockErr) StopUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, nil
}

func (m *dbusMockErr) RestartUnit(name string, mode string, ch chan<- string) (int, error) {
	return 0, nil
}

func TestMockGameBackendListGames(t *testing.T) {
	mock := &MockGameBackend{
		ListGamesFunc: func() ([]string, error) {
			return []string{"csgo", "rust"}, nil
		},
	}

	games, err := mock.ListGames()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 2 {
		t.Errorf("expected 2 games, got %d", len(games))
	}
	if games[0] != "csgo" {
		t.Errorf("expected first game 'csgo', got '%s'", games[0])
	}
}

func TestMockGameBackendListGamesError(t *testing.T) {
	mock := &MockGameBackend{
		ListGamesFunc: func() ([]string, error) {
			return nil, errors.New("dbus failed")
		},
	}

	_, err := mock.ListGames()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockGameBackendStartGame(t *testing.T) {
	called := false
	mock := &MockGameBackend{
		StartGameFunc: func(name string) error {
			called = true
			if name != "csgo" {
				t.Errorf("expected name 'csgo', got '%s'", name)
			}
			return nil
		},
	}

	err := mock.StartGame("csgo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected StartGameFunc to be called")
	}
}

func TestMockGameBackendStartGameError(t *testing.T) {
	mock := &MockGameBackend{
		StartGameFunc: func(name string) error {
			return errors.New("unit not found")
		},
	}

	err := mock.StartGame("missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockGameBackendStopGame(t *testing.T) {
	called := false
	mock := &MockGameBackend{
		StopGameFunc: func(name string) error {
			called = true
			return nil
		},
	}

	err := mock.StopGame("rust")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected StopGameFunc to be called")
	}
}

func TestMockGameBackendStopGameError(t *testing.T) {
	mock := &MockGameBackend{
		StopGameFunc: func(name string) error {
			return errors.New("permission denied")
		},
	}

	err := mock.StopGame("rust")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockGameBackendRestartGame(t *testing.T) {
	called := false
	mock := &MockGameBackend{
		RestartGameFunc: func(name string) error {
			called = true
			if name != "terraria" {
				t.Errorf("expected name 'terraria', got '%s'", name)
			}
			return nil
		},
	}

	err := mock.RestartGame("terraria")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected RestartGameFunc to be called")
	}
}

func TestMockGameBackendRestartGameError(t *testing.T) {
	mock := &MockGameBackend{
		RestartGameFunc: func(name string) error {
			return errors.New("service not running")
		},
	}

	err := mock.RestartGame("terraria")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockGameBackendGetGameStatus(t *testing.T) {
	mock := &MockGameBackend{
		GetGameStatusFunc: func(name string) (string, error) {
			if name != "csgo" {
				t.Errorf("expected name 'csgo', got '%s'", name)
			}
			return "active", nil
		},
	}

	status, err := mock.GetGameStatus("csgo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != "active" {
		t.Errorf("expected status 'active', got '%s'", status)
	}
}

func TestMockGameBackendGetGameStatusError(t *testing.T) {
	mock := &MockGameBackend{
		GetGameStatusFunc: func(name string) (string, error) {
			return "", errors.New("dbus error")
		},
	}

	_, err := mock.GetGameStatus("csgo")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMockGameBackendStreamLogs(t *testing.T) {
	var logs []string
	mock := &MockGameBackend{
		StreamLogsFunc: func(ctx context.Context, name string, callback func(string)) error {
			if name != "csgo" {
				t.Errorf("expected name 'csgo', got '%s'", name)
			}
			callback("log line 1")
			callback("log line 2")
			return nil
		},
	}

	err := mock.StreamLogs(context.Background(), "csgo", func(line string) {
		logs = append(logs, line)
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 log lines, got %d", len(logs))
	}
	if logs[0] != "log line 1" {
		t.Errorf("expected 'log line 1', got '%s'", logs[0])
	}
	if logs[1] != "log line 2" {
		t.Errorf("expected 'log line 2', got '%s'", logs[1])
	}
}

func TestMockGameBackendStreamLogsError(t *testing.T) {
	mock := &MockGameBackend{
		StreamLogsFunc: func(ctx context.Context, name string, callback func(string)) error {
			return errors.New("journal open failed")
		},
	}

	err := mock.StreamLogs(context.Background(), "csgo", func(line string) {})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
