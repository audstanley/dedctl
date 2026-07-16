package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"dedctl/internal/config"
)

// GameInfo holds enriched information about a game server.
type GameInfo struct {
	Name      string `json:"name"`
	AppId     int    `json:"app_id"`
	Order     int    `json:"order"`
	HasImage  bool   `json:"has_image"`
	Enabled   bool   `json:"enabled"`
	MainImage string `json:"main_image"`
	Icon      string `json:"icon"`
}

// ServerInfo holds global server metadata (public, no auth).
type ServerInfo struct {
	MainImage string `json:"main_image"`
	Icon      string `json:"icon"`
}

// GameBackend defines the interface for game server operations.
// Implementations can use systemd D-Bus, mock backends for testing, etc.
type GameBackend interface {
	ListGames() ([]string, error)
	ListGamesWithMeta() ([]GameInfo, error)
	StartGame(name string) error
	StopGame(name string) error
	RestartGame(name string) error
	EnableGame(name string) error
	DisableGame(name string) error
	GetGameStatus(name string) (string, error)
	StreamLogs(ctx context.Context, name string, callback func(string)) error
	UpdateMetadata(name string, appId, order int) error
	UpdateGlobalMetadata(field string, value string) error
	UpdateArt(name string, appId int) error
	GetServerInfo() ServerInfo
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

// Package-level metadata and image service instances set by the server startup.
var gameMetadata map[string]struct{ AppId int; Order int }
var imageService *ImageService
var imgDir string
var metaDir string
var gameMetadataObj *config.Metadata

// Default systemd user path for scanning service files
var systemdUserPath = filepath.Join(os.Getenv("HOME"), ".config", "systemd", "user")
var systemdWantsPath = filepath.Join(os.Getenv("HOME"), ".config", "systemd", "user", "default.target.wants")

// SetMetadataAndImages sets the metadata and image service references for use by game operations.
func SetMetadataAndImages(meta map[string]struct{ AppId int; Order int }, svc *ImageService, dir string) {
	gameMetadata = meta
	imageService = svc
	imgDir = dir
}

// SetMetaDir sets the metadata directory and metadata object for persistence operations.
func SetMetaDir(dir string, obj *config.Metadata) {
	metaDir = dir
	gameMetadataObj = obj
}

// execSystemctl runs systemctl --user <args>
func execSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", append([]string{"--user"}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl %s: %s (%w)", strings.Join(args, " "), string(output), err)
	}
	return nil
}

// isServiceEnabled checks if a game service is enabled (symlinked in default.target.wants)
func isServiceEnabled(name string) bool {
	wantsFile := filepath.Join(systemdWantsPath, fmt.Sprintf("steam-%s.service", name))
	_, err := os.Lstat(wantsFile)
	return err == nil
}

// discoverGames scans ~/.config/systemd/user/ for all steam-*.service files
func discoverGames() ([]string, error) {
	files, err := os.ReadDir(systemdUserPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read systemd user directory: %w", err)
	}

	var games []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), "steam-") && strings.HasSuffix(f.Name(), ".service") {
			gameName := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "steam-"), ".service")
			games = append(games, gameName)
		}
	}
	return games, nil
}

// NewGameService creates a new GameService connected to systemd D-Bus
func NewGameService() *GameService {
	conn, err := dbus.NewUserConnection()
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



// ListGames returns all available Steam games from systemd user service files
func (s *GameService) ListGames() ([]string, error) {
	return discoverGames()
}

// ListGamesWithMeta returns games with metadata enrichment including enabled state
func (s *GameService) ListGamesWithMeta() ([]GameInfo, error) {
	games, err := discoverGames()
	if err != nil {
		return nil, err
	}

	var infos []GameInfo
	for _, name := range games {
		gm, exists := gameMetadata[name]
		appId := 0
		order := 0
		if exists {
			appId = gm.AppId
			order = gm.Order
		}
		infos = append(infos, GameInfo{
			Name:      name,
			AppId:     appId,
			Order:     order,
			HasImage:  appId > 0 && imageService.ImageExists(appId, imgDir),
			Enabled:   isServiceEnabled(name),
			MainImage: gameMetadataObj.GetMainImage(),
			Icon:      gameMetadataObj.GetIcon(),
		})
	}

	return infos, nil
}

// StartGame starts a Steam game server (does NOT enable it)
func (s *GameService) StartGame(gameName string) error {
	return execSystemctl("start", fmt.Sprintf("steam-%s.service", gameName))
}

// StopGame stops a Steam game server (does NOT disable it)
func (s *GameService) StopGame(gameName string) error {
	return execSystemctl("stop", fmt.Sprintf("steam-%s.service", gameName))
}

// RestartGame restarts a Steam game server (does NOT affect enabled state)
func (s *GameService) RestartGame(gameName string) error {
	return execSystemctl("restart", fmt.Sprintf("steam-%s.service", gameName))
}

// EnableGame enables a Steam game to start on login (creates symlink)
func (s *GameService) EnableGame(gameName string) error {
	return execSystemctl("enable", fmt.Sprintf("steam-%s.service", gameName))
}

// DisableGame disables a Steam game from starting on login (removes symlink)
func (s *GameService) DisableGame(gameName string) error {
	return execSystemctl("disable", fmt.Sprintf("steam-%s.service", gameName))
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
func (s *GameService) StreamLogs(ctx context.Context, gameName string, callback func(string)) error {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return err
	}
	defer journal.Close()

	unitName := fmt.Sprintf("steam-%s.service", gameName)
	journal.AddMatch("_SYSTEMD_USER_UNIT=" + unitName)

	const maxHistoryEntries = 200

	// Seek to the tail and read backwards for the last N entries
	journal.SeekTail()
	var historyEntries []string
	var entriesRead int
	for {
		n, err := journal.Previous()
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}

		if entriesRead >= maxHistoryEntries {
			break
		}

		entry, err := journal.GetEntry()
		if err != nil {
			continue
		}

		msg, ok := entry.Fields["MESSAGE"]
		if !ok {
			continue
		}

		entriesRead++
		loc, _ := time.LoadLocation("America/Los_Angeles")
		logLine := fmt.Sprintf("[%s] %s", time.Unix(0, int64(entry.RealtimeTimestamp)*1000).In(loc).Format(time.RFC3339), msg)
		historyEntries = append(historyEntries, logLine)
	}
	// Reverse so entries are in chronological order
	for i, j := 0, len(historyEntries)-1; i < j; i, j = i+1, j-1 {
		historyEntries[i], historyEntries[j] = historyEntries[j], historyEntries[i]
	}
	for _, entry := range historyEntries {
		callback(entry)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

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

				msg, ok := entry.Fields["MESSAGE"]
				if !ok {
					continue
				}

				loc, _ := time.LoadLocation("America/Los_Angeles")
				logLine := fmt.Sprintf("[%s] %s", time.Unix(0, int64(entry.RealtimeTimestamp)*1000).In(loc).Format(time.RFC3339), msg)
				callback(logLine)
			}
		case sdjournal.SD_JOURNAL_NOP:
			continue
		case sdjournal.SD_JOURNAL_INVALIDATE:
			continue
		default:
			if status < 0 {
				return fmt.Errorf("error in Wait: %d", status)
			}
		}
	}
}

// UpdateMetadata updates the AppID and/or order for a game and persists to metadata.yaml.
func (s *GameService) UpdateMetadata(name string, appId, order int) error {
	if gameMetadataObj == nil {
		return fmt.Errorf("metadata not initialized")
	}

	gm, exists := gameMetadataObj.Games[name]
	if !exists {
		gm = config.GameMetadata{}
	}
	gm.AppId = appId
	gm.Order = order
	gameMetadataObj.Games[name] = gm

	// Update the package-level map too
	gameMetadata[name] = struct{ AppId int; Order int }{AppId: gm.AppId, Order: gm.Order}

	return config.SaveMetadata(metaDir, gameMetadataObj)
}

// UpdateGlobalMetadata updates a global metadata field (main_image or icon).
func (s *GameService) UpdateGlobalMetadata(field, value string) error {
	if gameMetadataObj == nil {
		return fmt.Errorf("metadata not initialized")
	}

	switch field {
	case "main_image":
		gameMetadataObj.SetMainImage(value)
	case "icon":
		gameMetadataObj.SetIcon(value)
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	return config.SaveMetadata(metaDir, gameMetadataObj)
}

// GetServerInfo returns global server metadata.
func (s *GameService) GetServerInfo() ServerInfo {
	if gameMetadataObj == nil {
		return ServerInfo{}
	}
	return ServerInfo{
		MainImage: gameMetadataObj.GetMainImage(),
		Icon:      gameMetadataObj.GetIcon(),
	}
}

// UpdateArt downloads the game cover image for a single game.
func (s *GameService) UpdateArt(name string, appId int) error {
	if imageService == nil || appId <= 0 {
		return fmt.Errorf("image service not initialized or invalid AppID")
	}
	return imageService.DownloadGameImage(appId, imgDir)
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
	ListGamesFunc            func() ([]string, error)
	ListGamesWithMetaFunc    func() ([]GameInfo, error)
	StartGameFunc            func(name string) error
	StopGameFunc             func(name string) error
	RestartGameFunc          func(name string) error
	EnableGameFunc           func(name string) error
	DisableGameFunc          func(name string) error
	GetGameStatusFunc        func(name string) (string, error)
	StreamLogsFunc           func(ctx context.Context, name string, callback func(string)) error
	UpdateMetadataFunc       func(name string, appId, order int) error
	UpdateGlobalMetadataFunc func(field, value string) error
	UpdateArtFunc            func(name string, appId int) error
	GetServerInfoFunc        func() ServerInfo
}

// ListGames implements GameBackend
func (m *MockGameBackend) ListGames() ([]string, error) {
	if m.ListGamesFunc != nil {
		return m.ListGamesFunc()
	}
	return []string{}, nil
}

// ListGamesWithMeta implements GameBackend
func (m *MockGameBackend) ListGamesWithMeta() ([]GameInfo, error) {
	if m.ListGamesWithMetaFunc != nil {
		return m.ListGamesWithMetaFunc()
	}
	return []GameInfo{}, nil
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

// EnableGame implements GameBackend
func (m *MockGameBackend) EnableGame(name string) error {
	return m.EnableGameFunc(name)
}

// DisableGame implements GameBackend
func (m *MockGameBackend) DisableGame(name string) error {
	return m.DisableGameFunc(name)
}

// GetGameStatus implements GameBackend
func (m *MockGameBackend) GetGameStatus(name string) (string, error) {
	return m.GetGameStatusFunc(name)
}

// StreamLogs implements GameBackend
func (m *MockGameBackend) StreamLogs(ctx context.Context, name string, callback func(string)) error {
	return m.StreamLogsFunc(ctx, name, callback)
}

// UpdateMetadata implements GameBackend
func (m *MockGameBackend) UpdateMetadata(name string, appId, order int) error {
	return m.UpdateMetadataFunc(name, appId, order)
}

// UpdateGlobalMetadata implements GameBackend
func (m *MockGameBackend) UpdateGlobalMetadata(field, value string) error {
	return m.UpdateGlobalMetadataFunc(field, value)
}

// UpdateArt implements GameBackend
func (m *MockGameBackend) UpdateArt(name string, appId int) error {
	return m.UpdateArtFunc(name, appId)
}

// GetServerInfo implements GameBackend
func (m *MockGameBackend) GetServerInfo() ServerInfo {
	return m.GetServerInfoFunc()
}
