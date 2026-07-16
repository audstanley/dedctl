package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// SteamAppResponse represents the response from the Steam Web API.
type SteamAppResponse struct {
	Data *SteamAppData `json:"data"`
}

// SteamAppData represents game metadata from the Steam Web API.
type SteamAppData struct {
	Name         string `json:"name"`
	HeaderImage  string `json:"header_image"`
	Type         string `json:"type"`
	IsFree       bool   `json:"is_free"`
	ShortDesc    string `json:"short_description"`
	HeaderImageV5 string `json:"header_image_v5"`
}

// ImageService handles downloading and managing game cover images.
type ImageService struct {
	httpGet func(url string) (*http.Response, error)
}

// NewImageService creates a new ImageService.
func NewImageService() *ImageService {
	return &ImageService{
		httpGet: http.Get,
	}
}

// GameImageInfo holds information about a cached game image.
type GameImageInfo struct {
	HasImage  bool   `json:"has_image"`
	ImagePath string `json:"-"`
}

// ImageExists checks if a cached image exists for the given AppID.
func (s *ImageService) ImageExists(appId int, imgDir string) bool {
	path := filepath.Join(imgDir, fmt.Sprintf("%d.jpg", appId))
	_, err := os.Stat(path)
	return err == nil
}

// GetImagePath returns the expected path for a cached game image.
func (s *ImageService) GetImagePath(appId int, imgDir string) string {
	return filepath.Join(imgDir, fmt.Sprintf("%d.jpg", appId))
}

// DownloadGameImage downloads the game cover image from the Steam Web API.
// It fetches the app details via https://store.steampowered.com/api/appdetails,
// extracts the header_image URL, and downloads it as a JPG to img/{appId}.jpg.
func (s *ImageService) DownloadGameImage(appId int, imgDir string) error {
	if appId <= 0 {
		return fmt.Errorf("invalid app ID: %d", appId)
	}

	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %w", err)
	}

	// Fetch app details from Steam Web API
	apiURL := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%d", appId)
	resp, err := s.httpGet(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam API for app %d: %w", appId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Steam API returned status %d for app %d", resp.StatusCode, appId)
	}

	var apiResp map[string]*SteamAppResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Steam API response: %w", err)
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse Steam API response: %w", err)
	}

	appData, ok := apiResp[fmt.Sprintf("%d", appId)]
	if !ok || appData == nil || appData.Data == nil {
		return fmt.Errorf("app %d not found on Steam", appId)
	}

	data := appData.Data
	if data.HeaderImage == "" {
		return fmt.Errorf("no header image available for app %d", appId)
	}

	imageURL := data.HeaderImage

	// Download the image
	imgResp, err := s.httpGet(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image for app %d: %w", appId, err)
	}
	defer imgResp.Body.Close()

	if imgResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", imgResp.StatusCode)
	}

	imgPath := s.GetImagePath(appId, imgDir)
	imgFile, err := os.Create(imgPath)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer imgFile.Close()

	if _, err := io.Copy(imgFile, imgResp.Body); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

// CacheMissingImages downloads images for all games that have an AppID set but no cached image.
func (s *ImageService) CacheMissingImages(games []string, meta map[string]struct{ AppId int; Order int }, imgDir string) error {
	for _, name := range games {
		gm, ok := meta[name]
		if !ok || gm.AppId <= 0 {
			continue
		}

		if s.ImageExists(gm.AppId, imgDir) {
			continue
		}

		if err := s.DownloadGameImage(gm.AppId, imgDir); err != nil {
			return fmt.Errorf("failed to cache image for %s (app %d): %w", name, gm.AppId, err)
		}
	}

	return nil
}

// CacheSingleImage downloads the image for a single game by name.
func (s *ImageService) CacheSingleImage(name string, appId int, imgDir string) error {
	return s.DownloadGameImage(appId, imgDir)
}
