package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMetadataNoFile(t *testing.T) {
	dir := t.TempDir()
	meta, err := LoadMetadata(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if meta.Games == nil {
		t.Error("expected Games map to be initialized")
	}
	if len(meta.Games) != 0 {
		t.Errorf("expected 0 games, got %d", len(meta.Games))
	}
}

func TestLoadMetadataWithFile(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `games:
  counter-strike:
    app_id: 730
    order: 1
  rust:
    app_id: 252490
    order: 2`

	if err := os.WriteFile(filepath.Join(dir, "metadata.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	meta, err := LoadMetadata(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(meta.Games) != 2 {
		t.Errorf("expected 2 games, got %d", len(meta.Games))
	}

	cs, ok := meta.Games["counter-strike"]
	if !ok {
		t.Error("expected counter-strike game")
	} else {
		if cs.AppId != 730 {
			t.Errorf("expected app_id 730, got %d", cs.AppId)
		}
		if cs.Order != 1 {
			t.Errorf("expected order 1, got %d", cs.Order)
		}
	}

	rust, ok := meta.Games["rust"]
	if !ok {
		t.Error("expected rust game")
	} else {
		if rust.AppId != 252490 {
			t.Errorf("expected app_id 252490, got %d", rust.AppId)
		}
	}
}

func TestSaveMetadata(t *testing.T) {
	dir := t.TempDir()
	meta := &Metadata{
		Games: map[string]GameMetadata{
			"csgo": {AppId: 730, Order: 1},
		},
	}

	if err := SaveMetadata(dir, meta); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadMetadata(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded.Games) != 1 {
		t.Errorf("expected 1 game, got %d", len(loaded.Games))
	}

	gm, ok := loaded.Games["csgo"]
	if !ok {
		t.Fatal("expected csgo game")
	}
	if gm.AppId != 730 {
		t.Errorf("expected app_id 730, got %d", gm.AppId)
	}
}

func TestSortedGames(t *testing.T) {
	meta := &Metadata{
		Games: map[string]GameMetadata{
			"zeus":        {Order: 0},
			"counter-strike": {Order: 1},
			"ark":         {Order: 0},
			"rust":        {Order: 2},
		},
	}

	games := meta.SortedGames()
	if len(games) != 4 {
		t.Fatalf("expected 4 games, got %d", len(games))
	}

	if games[0] != "counter-strike" {
		t.Errorf("expected counter-strike first, got %s", games[0])
	}
	if games[1] != "rust" {
		t.Errorf("expected rust second, got %s", games[1])
	}
	// Unordered games should be last, sorted alphabetically
	if games[2] != "ark" {
		t.Errorf("expected ark third (alphabetical), got %s", games[2])
	}
	if games[3] != "zeus" {
		t.Errorf("expected zeus fourth (alphabetical), got %s", games[3])
	}
}

func TestAddGame(t *testing.T) {
	meta := &Metadata{Games: map[string]GameMetadata{
		"rust": {Order: 1},
	}}

	meta.AddGame("csgo")
	if !meta.HasGame("csgo") {
		t.Error("expected csgo to be added")
	}
	if meta.Games["csgo"].Order != 2 {
		t.Errorf("expected order 2, got %d", meta.Games["csgo"].Order)
	}

	// Adding again should be a no-op
	meta.AddGame("csgo")
	if len(meta.Games) != 2 {
		t.Errorf("expected 2 games after duplicate add, got %d", len(meta.Games))
	}
}

func TestRemoveGame(t *testing.T) {
	meta := &Metadata{Games: map[string]GameMetadata{
		"rust": {Order: 1},
		"csgo": {Order: 2},
	}}

	meta.RemoveGame("csgo")
	if meta.HasGame("csgo") {
		t.Error("expected csgo to be removed")
	}
	if len(meta.Games) != 1 {
		t.Errorf("expected 1 game after remove, got %d", len(meta.Games))
	}
}

func TestGetAppId(t *testing.T) {
	meta := &Metadata{Games: map[string]GameMetadata{
		"csgo": {AppId: 730},
	}}

	if meta.GetAppId("csgo") != 730 {
		t.Errorf("expected app_id 730, got %d", meta.GetAppId("csgo"))
	}
	if meta.GetAppId("unknown") != 0 {
		t.Errorf("expected app_id 0 for unknown game, got %d", meta.GetAppId("unknown"))
	}
}

func TestSetAppId(t *testing.T) {
	meta := &Metadata{Games: map[string]GameMetadata{
		"csgo": {AppId: 0},
	}}

	meta.SetAppId("csgo", 730)
	if meta.GetAppId("csgo") != 730 {
		t.Errorf("expected app_id 730, got %d", meta.GetAppId("csgo"))
	}
}

func TestSetAppIdNewGame(t *testing.T) {
	meta := &Metadata{Games: map[string]GameMetadata{}}

	meta.SetAppId("new-game", 440)
	if meta.GetAppId("new-game") != 440 {
		t.Errorf("expected app_id 440, got %d", meta.GetAppId("new-game"))
	}
	if meta.Games["new-game"].Order != 0 {
		t.Errorf("expected order 0 for new game, got %d", meta.Games["new-game"].Order)
	}
}
