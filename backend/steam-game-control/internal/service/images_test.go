package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestImageExistsFalse(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	if svc.ImageExists(730, dir) {
		t.Error("expected image not to exist")
	}
}

func TestImageExistsTrue(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()

	path := svc.GetImagePath(730, dir)
	if err := os.WriteFile(path, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	if !svc.ImageExists(730, dir) {
		t.Error("expected image to exist")
	}
}

func TestGetImagePath(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	path := svc.GetImagePath(12345, dir)
	expected := dir + "/12345.jpg"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestImageExistsInvalidAppId(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	if svc.ImageExists(0, dir) {
		t.Error("expected image not to exist for app_id 0")
	}
}

func TestDownloadGameImageInvalidAppId(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	err := svc.DownloadGameImage(0, dir)
	if err == nil {
		t.Error("expected error for invalid app ID")
	}
}

func TestDownloadGameImageNotFound(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	err := svc.DownloadGameImage(999999999, dir)
	if err == nil {
		t.Error("expected error for non-existent app")
	}
}

func TestCacheMissingImagesSkipsNoAppId(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()

	meta := map[string]struct{ AppId int; Order int }{
		"game1": {AppId: 0, Order: 1},
	}

	err := svc.CacheMissingImages([]string{"game1"}, meta, dir)
	if err != nil {
		t.Errorf("expected no error for games without app_id, got %v", err)
	}
}

func TestCacheMissingImagesSkipsCached(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()

	path := svc.GetImagePath(730, dir)
	os.WriteFile(path, []byte("cached"), 0644)

	meta := map[string]struct{ AppId int; Order int }{
		"csgo": {AppId: 730, Order: 1},
	}

	err := svc.CacheMissingImages([]string{"csgo"}, meta, dir)
	if err != nil {
		t.Errorf("expected no error for cached image, got %v", err)
	}
}

func TestCacheSingleImageInvalidAppId(t *testing.T) {
	svc := NewImageService()
	dir := t.TempDir()
	err := svc.CacheSingleImage("game", 0, dir)
	if err == nil {
		t.Error("expected error for invalid app ID")
	}
}

func TestImageServiceCacheMissingImagesWithServer(t *testing.T) {
	svc := NewImageService()
	svc.httpGet = func(url string) (*http.Response, error) {
		if url == fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=730") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"730":{"success":true,"data":{"header_image":"https://example.com/730.jpg"}}}`))),
			}, nil
		}
		if url == "https://example.com/730.jpg" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte("test-image-data"))),
			}, nil
		}
		return nil, os.ErrNotExist
	}

	dir := t.TempDir()

	err := svc.DownloadGameImage(730, dir)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	path := svc.GetImagePath(730, dir)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved image: %v", err)
	}
	if string(data) != "test-image-data" {
		t.Errorf("expected test-image-data, got %s", string(data))
	}
}
