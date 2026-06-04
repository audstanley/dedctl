package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// GameMetadata holds metadata for a single game.
type GameMetadata struct {
	AppId int `yaml:"app_id,omitempty"`
	Order int `yaml:"order"`
}

// Metadata holds the game metadata configuration.
type Metadata struct {
	Games map[string]GameMetadata `yaml:"games"`
}

// LoadMetadata reads the metadata YAML file from the given directory.
func LoadMetadata(metaDir string) (*Metadata, error) {
	path := filepath.Join(metaDir, "metadata.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{Games: make(map[string]GameMetadata)}, nil
		}
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var meta Metadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	if meta.Games == nil {
		meta.Games = make(map[string]GameMetadata)
	}

	return &meta, nil
}

// SaveMetadata writes the metadata YAML file to the given directory.
func SaveMetadata(metaDir string, meta *Metadata) error {
	if meta.Games == nil {
		meta.Games = make(map[string]GameMetadata)
	}

	path := filepath.Join(metaDir, "metadata.yaml")

	data, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// SortedGames returns the games sorted by order, then alphabetically by name.
// Games without an order are sorted alphabetically at the end.
func (m *Metadata) SortedGames() []string {
	ordered := []string{}
	unordered := []string{}

	for name := range m.Games {
		if m.Games[name].Order > 0 {
			ordered = append(ordered, name)
		} else {
			unordered = append(unordered, name)
		}
	}

	sort.Strings(ordered)
	sort.Strings(unordered)

	return append(ordered, unordered...)
}

// HasGame returns true if the game exists in the metadata.
func (m *Metadata) HasGame(name string) bool {
	_, ok := m.Games[name]
	return ok
}

// AddGame adds a game to the metadata with the next available order.
func (m *Metadata) AddGame(name string) {
	if m.HasGame(name) {
		return
	}

	maxOrder := 0
	for _, gm := range m.Games {
		if gm.Order > maxOrder {
			maxOrder = gm.Order
		}
	}

	m.Games[name] = GameMetadata{Order: maxOrder + 1}
}

// RemoveGame removes a game from the metadata.
func (m *Metadata) RemoveGame(name string) {
	delete(m.Games, name)
}

// GetAppId returns the Steam AppID for a game, or 0 if not set.
func (m *Metadata) GetAppId(name string) int {
	gm, ok := m.Games[name]
	if !ok {
		return 0
	}
	return gm.AppId
}

// SetAppId sets the Steam AppID for a game.
func (m *Metadata) SetAppId(name string, appId int) {
	if m.Games == nil {
		m.Games = make(map[string]GameMetadata)
	}
	gm, ok := m.Games[name]
	if !ok {
		gm = GameMetadata{}
	}
	gm.AppId = appId
	m.Games[name] = gm
}
